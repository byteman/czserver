var app = angular.module('myApp')

app.controller('UpdateCtrl',['$scope', 'FileUploader', function($scope, FileUploader) {
    $scope.uploadStatus = $scope.uploadStatus1 = false; //定义两个上传后返回的状态，成功获失败
    var uploader = $scope.uploader = new FileUploader({
        url: '/upload',
        queueLimit: 1,     //文件个数 
        removeAfterUpload: true   //上传后删除文件
    });
    
    $scope.clearItems = function(){    //重新选择文件时，清空队列，达到覆盖文件的效果
        uploader.clearQueue();
    }
    
    uploader.onAfterAddingFile = function(fileItem) {
    	console.log("onAfterAddingFile")
        $scope.fileItem = fileItem._file;    //添加文件之后，把文件信息赋给scope
    };
   
    uploader.onSuccessItem = function(fileItem, response, status, headers) {
    	console.log("success")
        $scope.uploadStatus = true;   //上传成功则把状态改为true
    };
    
    $scope.UploadFile = function(){
        uploader.uploadAll();
       
    }
}])