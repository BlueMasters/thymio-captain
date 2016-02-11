/*
 * @author   Lucy Linder <lucy.derlin@gmail.com>
 * @date     Summer 2015
 *
 * Copyright 2015 Derlin. All rights reserved.
 * Use of this source code is governed by an Apache 2 license
 * that can be found in the LICENSE file.
 */
(function(){
    /**
     * @ngdoc controller
     * @name thymio.modals.DefaultModalController
     * @description
     * A default modal controller, which makes the object `attrs` available in the scope
     * and which calls `<caller>.close` upon close.
     *
     * See the example in the module's documentation page.
     */
    angular.module( 'thymio.modals' )
        .controller( 'DefaultModalController', DefaultModalController );

// --------------------------

    function DefaultModalController( $scope, attrs, close ){

        $scope.attrs = attrs;

        $scope.close = function( result ){
            close( result, 500 ); // close, but give 500ms for bootstrap to animate
        };

    }
}());
