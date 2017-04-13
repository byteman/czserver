var app = angular.module('myApp')

app.controller('weightCtrl', function($scope, $http,$routeParams) {
    
	// console.log("id=",$routeParams.id)


    $scope.getPageData = function(page){

        var url = "/weight?pages="+page;
        if($routeParams.id!=null)
        {
            url += "&id="+$routeParams.id;
        }
        console.log("data url=",url);
        $http.get(url)
            .success(function (response) {
                BootstrapPagination($("#pagination"), {
                    //记录总数。
                    total: response.Total,
                    //当前页索引编号。从其开始（从0开始）的整数。
                    pageIndex: page,
                    pageSize: 20,
                    pageSizeListFormateString:"",
                    //当分页更改后引发此事件。
                    pageChanged: function (pageIndex, pageSize) {

                        $scope.getPageData(pageIndex)
                    },
                });
                $scope.data = response.Weights;
            });
    }
	$scope.getPageData(0)
	
});