var app = angular.module('myApp')

app.controller('UpdateCtrl',['$scope', 'FileUploader', function($scope, FileUploader) {
    $scope.uploadStatus = $scope.uploadStatus1 = false; //���������ϴ��󷵻ص�״̬���ɹ���ʧ��
    var uploader = $scope.uploader = new FileUploader({
        url: '/upload',
        queueLimit: 1,     //�ļ����� 
        removeAfterUpload: true   //�ϴ���ɾ���ļ�
    });
    
    $scope.clearItems = function(){    //����ѡ���ļ�ʱ����ն��У��ﵽ�����ļ���Ч��
        uploader.clearQueue();
    }
    
    uploader.onAfterAddingFile = function(fileItem) {
    	console.log("onAfterAddingFile")
        $scope.fileItem = fileItem._file;    //����ļ�֮�󣬰��ļ���Ϣ����scope
    };
   
    uploader.onSuccessItem = function(fileItem, response, status, headers) {
    	console.log("success")
        $scope.uploadStatus = true;   //�ϴ��ɹ����״̬��Ϊtrue
    };
    
    $scope.UploadFile = function(){
        uploader.uploadAll();
       
    }
}])