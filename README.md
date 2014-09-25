cd simple-sample
   goapp deploy default/app.yaml
   goapp deploy mobile-frontend/mobile-frontend.yaml
   goapp deploy static-backend/my-module.yaml

Once the application has been successfully deployed 
you can access it at http://simple-sample.appspot.com. 
You can also access each of the modules individually:
   http://default.simple-sample.appspot.com
   http://mobile-frontend.simple-sample.appspot.com
   http://my-module.simple-sample.appspot.com
   
See https://developers.google.com/appengine/docs/go/modules/