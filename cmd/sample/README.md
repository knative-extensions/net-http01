# net-http01 sample

This small sample can be used to test the core libraries without all of the rest
of knative.dev/serving.

Instructions:

1. Deploy the load balance service.

   ```shell
   kubectl apply -f cmd/sample/service.yaml
   ```

1. Wait for the service to receive an IP or CNAME.

   ```shell
   watch kubectl get svc sample
   ```

1. Configure a DNS record for the IP or Hostname.

   If the service has an external IP, create an `A` record.

   If the service has a hostname, create a `CNAME` record.

1. Wait for the DNS record to go live.

   ```shell
   watch dig your-domain-name.io
   ```

1. Edit `sample.yaml` to pass the domain name.

   ```yaml
   args:
     - "-domain=your-domain-name.io"
   ```

1. Deploy the application.

   ```shell
   ko apply -Bf cmd/sample/sample.yaml
   ```

1. Curl the application:

   ```shell
   watch curl https://your-domain-name.io
   ```
