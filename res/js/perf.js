/* javascript */

var maxDataCount = 30;
var name = "fps"



// 基于准备好的dom，初始化echarts实例
var chartFps = echarts.init(document.getElementById('chart-fps'));
var chartSingleFps = echarts.init(document.getElementById('chart-single-fps'));

// 指定图表的配置项和数据

var fpsData = [];
for (var i = maxDataCount; i > 0; i -= 1) {
fpsData.push({
  value: [new Date().getTime() - 1000 * i, 0]
})
}

var fpsSingleData = [];


var option = {
title: {
  text: 'FPS'
},
toolbox: {
  feature: {
    saveAsImage: {}
  }
},
tooltip: {
  trigger: 'axis',
  // formatter: function(params) {
  // params = params[0];
  // console.log(params)
  // var date = new Date(params.value[0]);
  // return date + date.getFullYear() + '/' + (date.getMonth() + 1) + '/' + date.getDate() + ' : ' + params.value[1];
  // },
  axisPointer: {
    animation: false
  }
},
legend: {
  data: ['FPS']
},
xAxis: {
  type: 'time',
  splitLine: {
    show: false
  }
},
yAxis: {
  type: 'value',
  min: 0,
  max: 65,
  axisLabel: {
    formatter: '{value}'
  },
},
series: [{
  name: 'FPS',
  type: 'line',
  data: fpsData,
  animation: false,
  smooth: true,
  areaStyle: {
    normal: {}
  },
}]
}

chartFps.setOption(option);




var optionSingle = {
title: {
  text: 'FPS'
},
toolbox: {
  feature: {
    saveAsImage: {}
  }
},
tooltip: {
  trigger: 'axis',
  // formatter: function(params) {
  // params = params[0];
  // console.log(params)
  // var date = new Date(params.value[0]);
  // return date + date.getFullYear() + '/' + (date.getMonth() + 1) + '/' + date.getDate() + ' : ' + params.value[1];
  // },
  axisPointer: {
    animation: false
  }
},
legend: {
  data: ['FPS']
},
xAxis: {
  type: 'value',
  splitLine: {
    show: false
  }
},
yAxis: {
  type: 'value',
  min: 0,
  max: 65,
  axisLabel: {
    formatter: '{value}'
  },
},
series: [{
  name: 'FPS',
  type: 'line',
  data: fpsSingleData,
  animation: false,
  smooth: true,
  areaStyle: {
    normal: {}
  },
}]
}


var ws = newWebsocket('/ws/perfs/' + name, {
    onopen: function(evt) {
        console.log(evt);
    },
    onmessage: function(evt) {
        var data = JSON.parse(evt.data);
        if (fpsData && data.fps) {
            fpsData.push({
                value: [new Date(), data.fps],
            })
            if (fpsData.length > maxDataCount) {
                fpsData.shift();
            }
            chartFps.setOption({
                series: [{
                    data: fpsData,
                }]
            });
        }
    }
})




//add listener for btn

$("#start-test").click(function(){
  $.ajax({
    type:"get",
    url: "/start_fps",
    dataType: "json",
    success: function(data){
      console.log(data);
      console.log("start test.......");
      $("#start-test").prop("disabled", true);
      $("#end-test").prop("disabled", false);
    },
    error: function(){
        alert("无法开启帧率采集");
    }
  });

});


$("#end-test").click(function(){
  $("#modal-filename").modal("show");
});

$("#confirm-filename").click(function(){

    var filename = $("#input-filename").val();
    console.log(filename);
    $("#modal-filename").modal("hide");
    $.ajax({
    type:"get",
    url: "/stop_fps/" + filename,
    dataType: "json",
    success: function(data){
      console.log(data);
      console.log("end test.......");
      getFileList();
      $("#start-test").prop("disabled", false);
      $("#end-test").prop("disabled", true);
    },
    error: function(){
        alert("无法停止帧率采集");
    }
  });
});



getFileList();



function getFileList(){
  $.ajax({
    type:"get",
    url: "/api/file_list",
    dataType: "json",
    success: function(data){
      if(data.status == 1){
        $("#file-list").html("");
        for(var idx=0; idx<data.data.length; idx++){
          str = '<a class="list-group-item text-center" onclick="showFps(' + idx + ')" id="scene_' + idx+ '">' + data.data[idx] + '</a>'
          var existedContent = $("#file-list").html();
          $("#file-list").html(existedContent + str);
        }
      }
    },
    error: function(){
        alert("无法获取文件列表");
    }
  });
}


function showFps(idx) {
  var filename = $("#scene_" + idx).html();
  $.ajax({
    type:"get",
    url: "/api/fps_data/" + filename,
    dataType: "json",
    success: function(data){
      console.log(data);
      fpsSingleData = [];
      var maxFps = 0;
      var minFps = 100;
      var averFps = 0;
      for(var idx=0; idx< data.data.length; idx++){
        if(data.status == 1){
          fpsdata = parseInt(data.data[idx])
          if(fpsdata > maxFps){
            maxFps = fpsdata;
          }
          if(fpsdata < minFps){
            minFps = fpsdata;
          }
          averFps += fpsdata;
          fpsSingleData.push({
            value: [idx, fpsdata],
          })
        }
      }
      averFps = averFps / data.data.length;
      averFps = averFps.toFixed(2);

      chartSingleFps.setOption(optionSingle);
      chartSingleFps.setOption({
                series: [{
                    data: fpsSingleData,
                }]
            });
      $("#fps-title").html("场景 " + filename + " 帧率数据");
      $("#max-fps-data").html("最高帧率: " + maxFps);
      $("#min-fps-data").html("最低帧率: " + minFps);
      $("#aver-fps-data").html("平均帧率: " + averFps);
      $("#modal-fps").modal("show");
    },
    error: function(){
        alert("无法获取帧率");
    }
  });
}
