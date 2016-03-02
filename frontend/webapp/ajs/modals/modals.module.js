/*
 * @author   Lucy Linder <lucy.derlin@gmail.com>
 * @date     Summer 2015
 *
 * Copyright 2015 EIA-FR. All rights reserved.
 * Use of this source code is governed by an Apache 2 license
 * that can be found in the LICENSE file.
 */
(function(){
    /**
     * @ngdoc overview
     * @name derlin.modals
     * @description This module eases the process of showing
     * bootstrap-like modals from html templates.
     *
     * @example
     <pre class="prettyprint">
     <!--include bootstrap-->
     <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.0.0-alpha/css/bootstrap.css"
     type="text/css" rel="stylesheet">

     <!--body-->
     <div ng-controller="Ctrl">
     <button class="btn btn-default" ng-click="showModal()">show modal</button>
     <div>{{result}}</div>
     </div>


     <!--template: for yes/no modal with bootstrap -->
     <div class="modal fade">
     <div class="modal-dialog">
     <div class="modal-content">
     <div class="modal-header">
     <button type="button" class="close" ng-click="close(true)" data-dismiss="modal"
     aria-hidden="true">&times;</button>
     <h4 class="modal-title">Yes or No?</h4>
     </div>
     <div class="modal-body">
     <p>{{ attrs.text }}</p>
     </div>
     <div class="modal-footer">
     <button type="button" ng-click="close(false)" class="btn btn-default" data-dismiss="modal">No</button>
     <button type="button" ng-click="close(true)" class="btn btn-primary" data-dismiss="modal">Yes</button>
     </div>
     </div>
     </div>
     </div>

     <!--js-->
     <script>
     angular
     .module( 'modals.example', ['derlin.modals'] )
     .controller( 'Ctrl', function( $scope, ModalService ){
            $scope.showModal = function(){
                ModalService.showModal( {
                    templateUrl: 'YesNoModal.html',
                    controller : "DefaultModalController",
                    inputs     : {
                        attrs: {text: 'You have unsaved changes. Continue anyway ?'}
                    }

                } ).then( function( modal ){
                    modal.element.modal();
                    modal.close.then( function( result ){
                        $scope.result = result;
                    } );
                } );
            };
        } );
     </script>
     </pre>
     */
    angular.module( 'derlin.modals', [] );

}());