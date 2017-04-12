var app = angular.module('myApp')


app.controller('GpsListCtrl', function($scope, $http,$routeParams) {
    
	console.log("id=",$routeParams.id)
	$scope.isActivePage = function(page){
		return $scope.selPage == page;
	}
	
	$scope.getPageData = function(page){
		
		var url = "/gps?pages="+page;
		if($routeParams.id!=null)
		{
			url += "&id="+$routeParams.id;
		}
		console.log("url=",url);
		$http.get(url)
			.success(function (response) {
			$scope.data = response.gps;	
			$scope.total= response.Total
			$scope.pagesize= response.PageSize
			$scope.selectPage(page)
			console.log("total="+$scope.total+"size="+$scope.pagesize)
		});
	}
	
	//上一页
	$scope.Previous = function () {
		$scope.selectPage($scope.selPage - 1);
	}
	//下一页
	$scope.Next = function () {
		$scope.selectPage($scope.selPage + 1);
	}
	//打印当前选中页索引
	$scope.selectPage = function (page) {
		//不能小于1大于最大
		if (page < 1 || page > $scope.pagesize) return;
		//最多显示分页数5
		if (page > 2) {
			//因为只显示5个页数，大于2页开始分页转换
			var newpageList = [];
			for (var i = (page - 3) ; i < ((page + 2) > $scope.pages ? $scope.pages : (page + 2)) ; i++) {
			newpageList.push(i + 1);
		}
			$scope.pageList = newpageList;
		}
		$scope.selPage = page;
		//$scope.setData();
		$scope.isActivePage(page);
		//$scope.getPageData(page)
		console.log("选择的页：" + page);
	};
	
	$scope.getPageData(1)
	
});