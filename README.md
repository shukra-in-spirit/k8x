# k8x
Kube Chi - AI controlled Kubernetes scaling

Are you tired of having to regularly write promQL queries or pore over Grafana charts to set the scaling value for your kubernetes services? 
Or the slow and limited scaling offered by the Kubernetes Autoscaler is hurting your service latency?

Presenting k8x. The only Chi you need to find your peace.

Worry not. You can have latency and automation at the same time. Thanks to k8x. It offers planned automation. It uses Long Short Term Memory (LSTM) Neural Network to analyse the historical prometheus data for the service to be able to predict the correct size and replica for your service. It uses AWS Lambda and AWS S3 to generate, train, predict and retrain the LSTM models required to make the predictions for your service. The team can benefit from a hands free scaling of the kubernetes services that monitors regularly and scales about 20 minutes earlier to be able to meet the oncoming demand with the best possible manner.

How does k8x perform this magic?
Architecture
![image](https://github.com/shukra-in-spirit/k8x/assets/85339011/21079c13-37b3-4d09-a630-060279bf2bd1)
