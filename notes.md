## podman

## error in centos

```
cannot clone: Invalid argument user namespaces are not enabled in /proc/sys/user/max_user_namespaces Error: could not get runtime: cannot re-exec process podman

```
https://github.com/containers/podman/issues/7704

> `echo '63907' > /proc/sys/user/max_user_namespaces` this command solved it