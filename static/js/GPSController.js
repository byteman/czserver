var app = angular.module('myApp')

function addText(map,point,msg) {
    var opts = {
        position : point,    // 指定文本标注所在的地理位置
        offset   : new BMap.Size(30, -30)    //设置文本偏移量
    }
    var label = new BMap.Label(msg, opts);  // 创建文本标注对象
    label.setStyle({
        color : "red",
        fontSize : "12px",
        height : "20px",
        lineHeight : "20px",
        fontFamily:"微软雅黑"
    });
    map.addOverlay(label);
}
app.controller('GpsCtrl', function($scope, $http,$interval,$routeParams) {
	console.log("gps")
    var map = new BMap.Map("mymap");
	
	map.centerAndZoom(new BMap.Point(106.5478233651,29.5658377465),11);
    //map.setCurrentCity("上海");
	map.enableScrollWheelZoom(true);

	$scope.goPos=function(jd,wd){
		
		map.clearOverlays(); 
		var new_point = new BMap.Point(jd,wd);
		var marker = new BMap.Marker(new_point);  // 创建标注
		map.addOverlay(marker);              // 将标注添加到地图中
		map.panTo(new_point);      
	};
	 translateCallback = function (data){
      if(data.status === 0) {
          map.clearOverlays();
          var marker = new BMap.Marker(data.points[0]);
        map.addOverlay(marker);
        map.panTo(data.points[0]);
         // addText(map,data.points[0],"www")
      }
    }


	//$scope.goPos(106.5478233651,29.5658377465)
	$scope.getGps=function(page){
		
		var url = "/gps?pages="+page;
		if($routeParams.id!=null) {
            url += "&id=" + $routeParams.id;
        }
		$http.get(url).success(
		//$http.get("/gps?page=1&&id=3151").success(
			function (response) {
				$scope.gps = response.Gps;

				if($scope.gps.length>0)
				{

					var ggPoint = new BMap.Point($scope.gps[0].Longitude,$scope.gps[0].Latitude);
					var convertor = new BMap.Convertor();
					var pointArr = [];
					pointArr.push(ggPoint);
					convertor.translate(pointArr, 1, 5, translateCallback);



				}
				console.log($scope.gps)
			}
		);
	};
	var toDo = function () {
		
		$scope.getGps();
		
	};
	var timer = $interval(toDo, 1000);
	$scope.$on("$destroy",
		function( event ) {
			console.log("controller destroy,cancel timer")
			$interval.cancel( timer );

		}
    );
});