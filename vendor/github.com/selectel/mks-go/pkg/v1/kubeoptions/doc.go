/*
Package kubeoptions provides the ability to retrieve all available Kubernetes
feature gates and admission controllers through the MKS V1 API.

Example of getting available feature gates by Kubernetes version:

  availableFG, _, err := kubeoptions.ListFeatureGates(ctx, mksClient)
  if err != nil {
    log.Fatal(err)
  }
  for _, fgList := range availableFG {
    fmt.Printf("%s: %v\n", fgList.KubeVersion, fgList.Names)
  }

Example of getting available admission controllers by Kubernetes version:

  availableAC, _, err := kubeoptions.ListAdmissionControllers(ctx, mksClient)
  if err != nil {
    log.Fatal(err)
  }
  for _, acList := range availableAC {
    fmt.Printf("%s: %v\n", acList.KubeVersion, acList.Names)
  }
*/
package kubeoptions
