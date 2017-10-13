# Extensions for Gophercloud: an OpenStack SDK for Go

`gophercloud_extensions` extends [Gophercloud](https://github.com/gophercloud/gophercloud) by API requests to:

* blockstorage service list
* compute service list
* network agent list
* orchestration service list


The current main purpose is to use this data for monitoring. When this turns out to be useful, we might bring it back to Gophercloud.  The directory structure is borrowed from Gophercloud. 

## Testing

You need to Gophercloud to installed. Since this is a meant to be a library, we do not vendor it.

You might want to use `make test`
