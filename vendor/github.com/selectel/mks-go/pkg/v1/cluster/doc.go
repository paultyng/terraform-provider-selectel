/*
Package cluster provides the ability to retrieve and manage Kubernetes clusters
through the MKS V1 API.

Example of getting a single cluster referenced by its id

  mksCluster, _, err := cluster.Get(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", mksCluster)

Example of getting all clusters

  mksClusters, _, err := cluster.List(ctx, mksClient)
  if err != nil {
    log.Fatal(err)
  }
  for _, mksCluster := range mksClusters {
    fmt.Printf("%+v\n", mksCluster)
  }

Example of creating a new cluster

  createOpts := &cluster.CreateOpts{
    Name:        "test-cluster-0",
    KubeVersion: "1.15.7",
    Region:      "ru-1",
    Nodegroups: []*nodegroup.CreateOpts{
      {
        Count:            1,
        CPUs:             1,
        RAMMB:            2048,
        VolumeGB:         10,
        VolumeType:       "fast.ru-3a",
        KeypairName:      "ssh-key",
        AvailabilityZone: "ru-3a",
        Labels: map[string]string{
          "label-key0": "label-value0",
          "label-key1": "label-value1",
          "label-key2": "label-value2",
        },
        Taints: []nodegroup.Taint{
          {
            Key:    "test-key-0",
            Value:  "test-value-0",
            Effect: nodegroup.NoScheduleEffect,
          },
        },
      },
    },
  }
  mksCluster, _, err := cluster.Create(ctx, mksClient, createOpts)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", mksCluster)

Example of updating an existing cluster

  updateOpts := &cluster.UpdateOpts{
    MaintenanceWindowStart: "07:00:00",
    KubernetesOptions: &cluster.KubernetesOptions{
      EnablePodSecurityPolicy: false,
      FeatureGates: []string{
        "TTLAfterFinished",
      },
      AdmissionControllers: []string{
          "NamespaceLifecycle",
      },
    },
  }
  mksCluster, _, err := cluster.Update(ctx, mksClient, clusterID, updateOpts)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", mksCluster)

Example of deleting a single cluster

  _, err := cluster.Delete(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }

Example of getting a kubeconfig referenced by cluster id

  kubeconfig, _, err := cluster.GetKubeconfig(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Print(string(kubeconfig))

Example of rotating certificates by cluster id

  _, err := cluster.RotateCerts(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }

Example of upgrading Kubernetes patch version by cluster id

  mksCluster, _, err := cluster.UpgradePatchVersion(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", mksCluster)

Example of upgrading Kubernetes minor version by cluster id

  mksCluster, _, err := cluster.UpgradeMinorVersion(ctx, mksClient, clusterID)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%+v\n", mksCluster)
*/
package cluster
