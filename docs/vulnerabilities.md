
# Introduction

This is a brief tutorial on the use of `ion-connect` and the Ion Channel API to query the vulnerability database.

# Installation

The first step to using `ion-connect` is getting it installed and configured.
It can be installed from a prebuilt executable you can download the latest release here:

https://s3.amazonaws.com/public.ionchannel.io/files/ion-connect/ion-connect-latest.tar.gz

The package contains binaries compiled for linux, windows and mac os (darwin).

Once you get `ion-connect` properly installed and in your `PATH` you are ready to run the configure command.
The configure command will allow you to specify your api key (provided by Ion Channel) and any default properties.

```
$ ion-connect configure
Ion Channel Api Key []:
```

This will create a `$HOME/.ionchannel/credentials` file that will be used for authentication during the following commands.
NOTE: You can override the value in the `$HOME/.ionchannel/credentials` file at anytime using the `IONCHANNEL_SECRET_KEY` environment variable.

If your `ion-connect` environment hasn't been configured, or you don't have the env var set you'll see:

```
$ ion-connect vulnerability get-vulnerabilities
Unauthorized, make sure you run 'ion-connect configure' and set your Api Token
```

# Vulnerabilities

With a working `ion-connect` you're now ready to experiment with a few queries.  You can retrieve a list of actions, or commands, for a given resource - in this case we are interested in vulnerabilities.

```
➜  ~ ion-connect vulnerability
NAME:
   ion-connect vulnerability - set of commands for querying vulnerabilities

USAGE:
   ion-connect vulnerability command [command options] [arguments...]

COMMANDS:
   get-vulnerabilities-in-file    returns a list of vulnerabilities for a given set of product/version
   get-vulnerabilities-for-list    returns a list of vulnerabilities for a given set of product/version
   get-vulnerability        returns a specific vulnerability for a given query
   get-vulnerabilities        returns a list of vulnerabilities for a given query
   get-product            returns a specific product for a given query
   get-products            returns a list of products for a given query
   help, h            Shows a list of commands or help for one command

OPTIONS:
   --help, -h    show help
   ```

## get-products

Let's start with a simple query.  Ion Channel treats, pretty much everything, as a combination of organization or vendor and product or project - in the form of org/product.  This might be apache/tomcat or docker/docker or oracle/jdk.  Queries can be done without the organization prefix.  The results refer to CPEs in
the NIST NVD database.

Let's try docker/docker:

```
➜  ~ ion-connect vulnerability get-products docker/docker
[
  {
    "cpe23_name": "cpe:2.3:a:docker:docker:1.0.0:*:*:*:*:*:*:*",
    "created_at": "2016-08-19T16:34:13.137Z",
    "edition": "*",
    "id": 28794,
    "language": "*",
    "name": "cpe:/a:docker:docker:1.0.0",
    "other": "*",
    "part": "a",
    "product": "docker",
    "references": {
      "vendor product information": "https://www.docker.com/"
    },
    "sw_edition": "*",
    "target_hw": "*",
    "target_sw": "*",
    "title": "Docker 1.0.0",
    "update_id": "*",
    "updated_at": "2016-08-19T16:34:13.137Z",
    "vendor": "docker",
    "version": "1.0.0",
    "vulnerable_software_list": "cpe:/a:docker:docker:1.0.0"
  },
...
]
```

This returns a list of products that are related to the product: docker, for the organization docker.

## get-product

It is also possible to query for a specific CPE name:

```
➜  ~ ion-connect vulnerability get-product cpe:/a:docker:docker:1.0.0
{
  "cpe23_name": "cpe:2.3:a:docker:docker:1.0.0:*:*:*:*:*:*:*",
  "created_at": "2016-08-19T16:34:13.137Z",
  "edition": "*",
  "id": 28794,
  "language": "*",
  "name": "cpe:/a:docker:docker:1.0.0",
  "other": "*",
  "part": "a",
  "product": "docker",
  "references": {
    "vendor product information": "https://www.docker.com/"
  },
  "sw_edition": "*",
  "target_hw": "*",
  "target_sw": "*",
  "title": "Docker 1.0.0",
  "update_id": "*",
  "updated_at": "2016-08-19T16:34:13.137Z",
  "vendor": "docker",
  "version": "1.0.0",
  "vulnerable_software_list": "cpe:/a:docker:docker:1.0.0"
}
```

or the CPE in 2.3 form:

```
➜  ~ ion-connect vulnerability get-product 'cpe:2.3:a:docker:docker:1.0.0:*:*:*:*:*:*:*'
{
  "cpe23_name": "cpe:2.3:a:docker:docker:1.0.0:*:*:*:*:*:*:*",
  "created_at": "2016-08-19T16:34:13.137Z",
  "edition": "*",
  "id": 28794,
  "language": "*",
  "name": "cpe:/a:docker:docker:1.0.0",
  "other": "*",
  "part": "a",
  "product": "docker",
  "references": {
    "vendor product information": "https://www.docker.com/"
  },
  "sw_edition": "*",
  "target_hw": "*",
  "target_sw": "*",
  "title": "Docker 1.0.0",
  "update_id": "*",
  "updated_at": "2016-08-19T16:34:13.137Z",
  "vendor": "docker",
  "version": "1.0.0",
  "vulnerable_software_list": "cpe:/a:docker:docker:1.0.0"
}
```

## get-vulnerabilities

Much like with Products, vulnerabilities are queried by org/product.  But, for vulnerabilities you can also include a version.

```
➜  ~ ion-connect vulnerability get-vulnerabilities docker/docker 1.6
[
  {
    "access_complexity": "LOW",
    "access_vector": "LOCAL",
    "assessment_check": {},
    "availability_impact": "PARTIAL",
    "confidentiality_impact": "NONE",
    "created_at": "2016-08-19T18:19:04.516Z",
    "cve_name": "CVE-2015-3631",
    "date_generated": "2015-07-02T14:36:22.810-04:00",
    "date_published": "2015-05-18T11:59:17.760-04:00",
    "id": 71424,
    "integrity_impact": "PARTIAL",
    "last_modified": "2015-07-02T22:39:34.887-04:00",
    "references": [
      {
        "reference": "https://groups.google.com/forum/#!searchin/docker-user/1.6.1/docker-user/47GZrihtr-4/nwgeOOFLexIJ",
        "reference_type": "UNKNOWN",
        "source": "CONFIRM",
        "xml:lang": "en"
      },
      {
        "reference": "20150508 Docker 1.6.1 - Security Advisory [150507]",
        "reference_type": "UNKNOWN",
        "source": "FULLDISC",
        "xml:lang": "en"
      },
      {
        "reference": "http://packetstormsecurity.com/files/131835/Docker-Privilege-Escalation-Information-Disclosure.html",
        "reference_type": "UNKNOWN",
        "source": "MISC",
        "xml:lang": "en"
      }
    ],
    "scanner": {},
    "summary": "Docker Engine before 161 allows local users to set arbitrary Linux Security Modules LSM and docker_t policies via an image that allows volumes to override files in proc",
    "updated_at": "2016-08-19T18:19:04.516Z",
    "vulnerability_authentication": "NONE",
    "vulnerability_score": "3.6",
    "vulnerability_source": "http://nvd.nist.gov",
    "vulnerable_software_list": [
      {
        "cpe23_name": "cpe:2.3:a:docker:docker:1.6:*:*:*:*:*:*:*",
        "created_at": "2016-08-19T16:34:13.159Z",
        "edition": "*",
        "id": 28799,
        "language": "*",
        "name": "cpe:/a:docker:docker:1.6",
        "other": "*",
        "part": "a",
        "product": "docker",
        "references": {
          "Product Changelog": "https://docs.docker.com/release-notes/",
          "Vendor Website": "https://www.docker.com/"
        },
        "sw_edition": "*",
        "target_hw": "*",
        "target_sw": "*",
        "title": "Docker 1.6",
        "update_id": "*",
        "updated_at": "2016-08-19T16:34:13.159Z",
        "vendor": "docker",
        "version": "1.6",
        "vulnerable_software_list": "cpe:/a:docker:docker:1.6"
      }
    ]
  },
...
]
```

## get-vulnerability

And it is possible to query for a vulnerability's details by ID/name:

```
➜  ~ ion-connect vulnerability get-vulnerability CVE-2015-3631
{
  "access_complexity": "LOW",
  "access_vector": "LOCAL",
  "assessment_check": {},
  "availability_impact": "PARTIAL",
  "confidentiality_impact": "NONE",
  "cve_name": "CVE-2015-3631",
  "date_generated": "2015-07-02T14:36:22.810-04:00",
  "date_published": "2015-05-18T11:59:17.760-04:00",
  "integrity_impact": "PARTIAL",
  "last_modified": "2015-07-02T22:39:34.887-04:00",
  "references": [
    {
      "reference": "https://groups.google.com/forum/#!searchin/docker-user/1.6.1/docker-user/47GZrihtr-4/nwgeOOFLexIJ",
      "reference_type": "UNKNOWN",
      "source": "CONFIRM",
      "xml:lang": "en"
    },
    {
      "reference": "20150508 Docker 1.6.1 - Security Advisory [150507]",
      "reference_type": "UNKNOWN",
      "source": "FULLDISC",
      "xml:lang": "en"
    },
    {
      "reference": "http://packetstormsecurity.com/files/131835/Docker-Privilege-Escalation-Information-Disclosure.html",
      "reference_type": "UNKNOWN",
      "source": "MISC",
      "xml:lang": "en"
    }
  ],
  "scanner": {},
  "summary": "Docker Engine before 161 allows local users to set arbitrary Linux Security Modules LSM and docker_t policies via an image that allows volumes to override files in proc",
  "vulnerability_authentication": "NONE",
  "vulnerability_score": "3.6",
  "vulnerability_source": "http://nvd.nist.gov",
  "vulnerable_software_list": [
    {
      "cpe23_name": "cpe:2.3:a:docker:docker:1.6:*:*:*:*:*:*:*",
      "edition": "*",
      "language": "*",
      "name": "cpe:/a:docker:docker:1.6",
      "other": "*",
      "part": "a",
      "product": "docker",
      "references": {
        "Product Changelog": "https://docs.docker.com/release-notes/",
        "Vendor Website": "https://www.docker.com/"
      },
      "sw_edition": "*",
      "target_hw": "*",
      "target_sw": "*",
      "title": "Docker 1.6",
      "update_id": "*",
      "vendor": "docker",
      "version": "1.6"
    }
  ]
}
```

## get-vulnerabilities-for-file

It is possible to have a list inside a file, in the form of:

```
{"data":[{"product":"docker/docker","version":"1.6"},{"product":"redhat/docker"}]}
```

And reference it from a call:

```
➜  ~ ion-connect vulnerability get-vulnerabilities-in-file /somedir/somefile.json
[
  {
    "product": "docker/docker",
    "version": "1.6",
    "vulnerabilities": [
      {
        "access_complexity": "LOW",
        "access_vector": "LOCAL",
        "assessment_check": {},
...
```

## get-vulnerabilities-in-list

Like with the file, a list can be passed to get a single result set:

```
➜  ~ ion-connect vulnerability get-vulnerabilities-for-list '[{"product":"docker/docker","version":"1.6"},{"product":"redhat/docker"}]'
[
  {
    "product": "docker/docker",
    "version": "1.6",
    "vulnerabilities": [
      {
        "access_complexity": "LOW",
        "access_vector": "LOCAL",
        "assessment_check": {},
        "availability_impact": "PARTIAL",
        "confidentiality_impact": "NONE",
        "created_at": "2016-08-19T18:19:04.516Z",
        "cve_name": "CVE-2015-3631",
...
```


# HTTPS API

The `ion-connect` CLI is a basic wrapper around a handful of HTTP APIs that Ion Channel exposes through an API gateway.  The API endpoints are accessible directly with the same token that is used in the `ion-connect` configuration.

[NOTE: we use `jq` to help in processing the JSON responses.]

The full API specification is here: [http://ion.docs.ionchannel.io/bunsen.html](http://ion.docs.ionchannel.io/bunsen.html)

## getProducts?product=

```
curl https://api.ionchannel.io/v1/vulnerability/getVulnerabilities\?product\=docker/docker\&apikey\=3c92b12c2b0344a0834f21589d35df87
```

returns:

```
➜  ~ curl https://api.ionchannel.io/v1/vulnerability/getProducts\?product\=docker/docker\&apikey\=APIKEY | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  3739  100  3739    0     0   7402      0 --:--:-- --:--:-- --:--:--  7418
{
  "links": {
    "self": "https://api.ionchannel.io/v1/vulnerability/getProducts?apikey=APIKEY&product=docker%2fdocker",
    "next": null,
    "prev": null
  },
  "meta": {
    "copyright": "Copyright 2016, Ion Channel (ionchannel.io)",
    "authors": [
      "kellyp",
      "kitplummer"
    ],
    "version": "v1",
    "terms": "http://ionchannel.io/terms_privacy.html",
    "offset": 0,
    "limit": 10,
    "total_count": 6
  },
  "data": [
    {
      "id": 28794,
      "product": "docker",
      "version": "1.0.0",
      "vendor": "docker",
      "update_id": "*",
      "name": "cpe:/a:docker:docker:1.0.0",
      "part": "a",
      "edition": "*",
      "language": "*",
      "vulnerable_software_list": "cpe:/a:docker:docker:1.0.0",
      "references": {
        "vendor product information": "https://www.docker.com/"
      },
      "title": "Docker 1.0.0",
      "cpe23_name": "cpe:2.3:a:docker:docker:1.0.0:*:*:*:*:*:*:*",
      "sw_edition": "*",
      "target_sw": "*",
      "target_hw": "*",
      "other": "*",
      "created_at": "2016-08-19T16:34:13.137Z",
      "updated_at": "2016-08-19T16:34:13.137Z"
    },
...
```

## getProducts?name=

```
➜  ~ curl https://api.ionchannel.io/v1/vulnerability/getProducts\?name\=cpe:2.3:a:docker:docker:1.6:\*:\*:\*:\*:\*:\*:\*\&apikey\=APIKEY | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1051  100  1051    0     0   1600      0 --:--:-- --:--:-- --:--:--  1599
{
  "links": {
    "self": "https://api.ionchannel.io/v1/vulnerability/getProducts?apikey=APIKEY&name=cpe%3a2%2e3%3aa%3adocker%3adocker%3a1%2e6%3a%2a%3a%2a%3a%2a%3a%2a%3a%2a%3a%2a%3a%2a"
  },
  "meta": {
    "copyright": "Copyright 2016, Ion Channel (ionchannel.io)",
    "authors": [
      "kellyp",
      "kitplummer"
    ],
    "version": "v1",
    "terms": "http://ionchannel.io/terms_privacy.html",
    "offset": 0,
    "limit": 10
  },
  "data": {
    "id": 28799,
    "product": "docker",
    "version": "1.6",
    "vendor": "docker",
    "update_id": "*",
    "name": "cpe:/a:docker:docker:1.6",
    "part": "a",
    "edition": "*",
    "language": "*",
    "vulnerable_software_list": "cpe:/a:docker:docker:1.6",
    "references": {
      "Vendor Website": "https://www.docker.com/",
      "Product Changelog": "https://docs.docker.com/release-notes/"
    },
    "title": "Docker 1.6",
    "cpe23_name": "cpe:2.3:a:docker:docker:1.6:*:*:*:*:*:*:*",
    "sw_edition": "*",
    "target_sw": "*",
    "target_hw": "*",
    "other": "*",
    "created_at": "2016-08-19T16:34:13.159Z",
    "updated_at": "2016-08-19T16:34:13.159Z"
  },
  "timestamps": {
    "created": "2016-08-29T19:28:05.551+00:00",
    "updated": "2016-08-29T19:28:05.576+00:00"
  }
}
```

## getVulnerabilities

```
➜  ~ curl https://api.ionchannel.io/v1/vulnerability/getVulnerabilities\?product\=docker/docker\&apikey\=APIKEY | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 19656  100 19656    0     0  22608      0 --:--:-- --:--:-- --:--:-- 22593
{
  "links": {
    "self": "https://api.ionchannel.io/v1/vulnerability/getVulnerabilities?apikey=APIKEY&product=docker%2fdocker",
    "next": null,
    "prev": null
  },
  "meta": {
    "copyright": "Copyright 2016, Ion Channel (ionchannel.io)",
    "authors": [
      "kellyp",
      "kitplummer"
    ],
    "version": "v1",
    "terms": "http://ionchannel.io/terms_privacy.html",
    "offset": 0,
    "limit": 10,
    "total_count": 9
  },
  "data": [
    {
      "id": 71424,
      "cve_name": "CVE-2015-3631",
      "summary": "Docker Engine before 161 allows local users to set arbitrary Linux Security Modules LSM and docker_t policies via an image that allows volumes to override files in proc",
...
```

## getVulnerability

```
➜  ~ curl https://api.ionchannel.io/v1/vulnerability/getVulnerability\?name\=CVE-2015-3631\&apikey\=APIKEY | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  2015  100  2015    0     0   1985      0  0:00:01  0:00:01 --:--:--  1987
{
  "links": {
    "self": "https://api.ionchannel.io/v1/vulnerability/getVulnerability?apikey=APIKEY&name=CVE%2d2015%2d3631"
  },
  "meta": {
    "copyright": "Copyright 2016, Ion Channel (ionchannel.io)",
    "authors": [
      "kellyp",
      "kitplummer"
    ],
    "version": "v1",
    "terms": "http://ionchannel.io/terms_privacy.html",
    "offset": 0,
    "limit": 10
  },
  "data": {
    "cve_name": "CVE-2015-3631",
    "date_published": "2015-05-18T11:59:17.760-04:00",
    "last_modified": "2015-07-02T22:39:34.887-04:00",
...
```
