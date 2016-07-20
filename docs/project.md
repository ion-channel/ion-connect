# Introduction

The following document details the steps needed to use `ion-connect` to setup a project
and begin analysis.  With these steps you will be able to view current rules, create and view
rulesets and create, view, update and analyze a given project.  The analysis steps can easily be added to a CI/CD pipeline for monitoring risk and compliance automatically.


# Install

The first step when using `ion-connect` is getting it installed and configured.  
It can be installed from a prebuilt executable you can download for here (TODO: figure out where to host it).

Once you get `ion-connect` properly installed and in your `PATH` you are ready to run the `configure` command.  
The `configure` command will allow you to specify your api key (provided by Ion Channel along with your account id) and any default properties.


```
$ ion-connect configure
Ion Channel Api Key []:
```

This will create a `$HOME/.aws/credentials` file that will be used for authentication during the following commands.

NOTE:  You can override the value in the `$HOME/.aws/credentials` file at anytime using the `IONCHANNEL_SECRET_KEY` environment variable.

# Rulesets

Now you are ready to run ion-connect and set up some rule sets and projects to analyze for compliance and risk.  The first step is to view a list of rules provided by Ion Channel.

```
$ ion-connect ruleset get-rules
[
  {
    "category": "About Dot Yaml",
    "description": "The project source is required to include a valid .about.yml file.",
    "id": "c30b917956c3040daa2c571ef31dbe3a",
    "name": "Has a valid .about.yml file",
    "policy_url": "url",
    "remediation_url": "url",
    "scan_type": "about_yml"
  },
  ...
  {
    "category": "Code Coverage",
    "description": "A longer description of the rule: Code Coverage > 90%",
    "id": "865b64e5d9d936ced71582e88146dd11",
    "name": "Code Coverage > 90%",
    "policy_url": "url",
    "remediation_url": "url",
    "scan_type": "coverage"
  }
]
```

From here you can define rule sets which will govern your projects.  The following example will create a ruleset that requires a valid .about.yml file.

```
$ ion-connect ruleset create-ruleset --account-id <your-account-id> "rule set name" "this is a test ruleset" '[c30b917956c3040daa2c571ef31dbe3a"]'
```

At any time you can query for all or any of the rule sets in your account using the `get-rulesets` and `get-ruleset` commands.

```
$ ion-connect ion-connect ruleset get-rulesets --account-id <your-account-id>
[
  {
    "account_id": "<your-account-id>",
    "created_at": "2016-07-19T21:52:33.312+00:00",
    "description": "this is a test ruleset",
    "id": "<some-ruleset-id>",
    "name": "rule set name",
    "rule_ids": [
      "c30b917956c3040daa2c571ef31dbe3a"
    ],
    "rules": [
      {
        "category": "About Dot Yaml",
        "description": "The project source is required to include a valid .about.yml file.",
        "id": "c30b917956c3040daa2c571ef31dbe3a",
        "name": "Has a valid .about.yml file",
        "policy_url": "url",
        "remediation_url": "url",
        "scan_type": "about_yml"
      }
    ],
    "updated_at": "2016-07-19T21:52:33.312+00:00"
  }
]
```

and

```
$ ion-connect ruleset get-ruleset --account-id <your-account-id> <some-ruleset-id>
{
  "account_id": "<some-ruleset-id>",
  "created_at": "2016-07-19T21:52:33.312+00:00",
  "description": "this is a test ruleset",
  "id": "some-ruleset-id",
  "name": "rule set name",
  "rule_ids": [
    "c30b917956c3040daa2c571ef31dbe3a"
  ],
  "rules": [
    {
      "category": "About Dot Yaml",
      "description": "The project source is required to include a valid .about.yml file.",
      "id": "c30b917956c3040daa2c571ef31dbe3a",
      "name": "Has a valid .about.yml file",
      "policy_url": "url",
      "remediation_url": "url",
      "scan_type": "about_yml"
    }
  ],
  "updated_at": "2016-07-19T21:52:33.312+00:00"
}
```

# Projects

After you have a rule set defined you can create your project in ion channel for analysis.  The following will create a project record in Ion Channel named 'Project Name' for analysis of the DevOps/sonar-auth-geoaxis project in gitlab.  Since we are using the rule set previously defined, analysis of this project will fail unless there is a valid .about.yml file at the root of the repository.

```
$ ion-connect project create-project --account-id <your-account-id> --ruleset-id <some-ruleset-id>  --active "Project Name" "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git" "Project description"
{
  "account_id": "<your-account-id>",
  "active": true,
  "branch": "master",
  "created_at": "2016-07-19T22:27:23.646Z",
  "description": "Project description",
  "id": "<some-project-id>",
  "name": "Project Name",
  "ruleset_id": "<some-ruleset-id>",
  "source": "git@gitlab.devops.geointservices.io:GEOINTApps/MAGE.git",
  "type": "git",
  "updated_at": "2016-07-19T22:27:23.646Z"
}
```

Similar to rule sets, `ion-connect` provides commands for querying for the projects in your account.  

```
$ ion-connect project get-project --account-id <your-account-id> <some-project-id>
{
  "account_id": "<your-account-id>",
  "active": true,
  "branch": "master",
  "created_at": "2016-07-19T22:27:23.646Z",
  "description": "Project description",
  "id": "<some-project-id>",
  "name": "Project Name",
  "ruleset_id": "<some-ruleset-id>",
  "source": "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git",
  "type": "git",
  "updated_at": "2016-07-19T22:27:23.646Z"
}
```

# Analysis

Now that you have your project record in Ion Channel it's time to do some analysis. This is done with the `ion-connect scanner analyze-project` project command on the scanner resource.

NOTE: This is meant to provide an example of you to manually run the analysis of a project.  The following commands should be combined for use inside a CI/CD pipeline to ensure compliance is met.

```
ion-connect scanner analyze-project --account-id <your-account-id> --project-id <some-project-id> 1
{
  "account_id": "<your-account-id>",
  "build_number": "1",
  "created_at": "2016-07-19T22:46:29.123Z",
  "id": "<some-analysis-id>",
  "message": "Request for analysis  on Project Name has been accepted.",
  "project_id": "<some-project-id>",
  "status": "accepted",
  "updated_at": "2016-07-19T22:46:29.123Z"
}
```

Since the analysis happens asychronously you can monitor the status of the analysis with the `ion-connect scanner get-analysis-status` command.  Once the analysis has completed you should see some like this:

```
ion-connect scanner get-analysis-status --account-id <your-account-id> --project-id <some-project-id> <some-analysis-id>
{
  "account_id": "<your-account-id>",
  "build_number": "1",
  "created_at": "2016-07-20T16:43:49.386Z",
  "id": "<some-analysis-id>",
  "message": null,
  "project_id": "<some-project-id>",
  "scan_status": [
    {
      "account_id": "<your-account-id>",
      "analysis_status_id": "<some-analysis-id>",
      "created_at": "2016-07-20T16:43:50.078Z",
      "id": "784c6356-0508-6cd6-0ea6-3bd20d8268a5",
      "message": "Finished dependency scan for Sonar-Plugin, found  out of version dependencies.",
      "name": "dependency",
      "project_id": "<some-project-id>",
      "read": "f",
      "status": "finished",
      "updated_at": "2016-07-20T16:43:50.090Z"
    },
    {
      "account_id": "<your-account-id>",
      "analysis_status_id": "<some-analysis-id>",
      "created_at": "2016-07-20T16:43:50.347Z",
      "id": "640af301-cb23-b2e9-8614-7d74479ee23a",
      "message": "Finished file type scan for Sonar-Plugin, all file types look correct.",
      "name": "file_type",
      "project_id": "<some-project-id>",
      "read": "f",
      "status": "finished",
      "updated_at": "2016-07-20T16:43:50.358Z"
    },
    {
      "account_id": "<your-account-id>",
      "analysis_status_id": "<some-analysis-id>",
      "created_at": "2016-07-20T16:43:50.648Z",
      "id": "80c00035-3aff-cd42-a53d-cb42d64a5e95",
      "message": "Finished about_yml scan for Sonar-Plugin, valid .about.yml found.",
      "name": "about_yml",
      "project_id": "<some-project-id>",
      "read": "f",
      "status": "finished",
      "updated_at": "2016-07-20T16:43:50.652Z"
    },
    {
      "account_id": "<your-account-id>",
      "analysis_status_id": "<some-analysis-id>",
      "created_at": "2016-07-20T16:43:58.320Z",
      "id": "c91e9209-d5c8-764b-5f96-cc5130454e95",
      "message": "Finished clamav scan for Sonar-Plugin, found 0 infected files.",
      "name": "clamav",
      "project_id": "<some-project-id>",
      "read": "f",
      "status": "finished",
      "updated_at": "2016-07-20T16:43:58.324Z"
    }
  ],
  "status": "finished",
  "updated_at": "2016-07-20T16:43:58.379Z"
}
```


You can see from the above output the analysis finished with several scans also completing.  Once the analysis is finished you can request the evaluated analysis results with an additional command.  This details below show that the analysis completed

```
$ ion-connect analysis get-analysis --account-id <your-account-id> --project-id <some-project-id> <some-analysis-id>
{
  "account_id": "<your-account-id>",
  "branch": "master",
  "build_number": "1",
  "created_at": "2016-07-20T16:43:49.840Z",
  "description": "",
  "duration": 16305,
  "id": "<some-analysis-id>",
  "name": "compliance analysis",
  "passed": true,
  "project_id": "<some-project-id>",
  "risk": "low",
  "ruleset_id": "157b5aef071ab5d2ae3182b54da00f82",
  "ruleset_name": "rule set name",
  "scan_summaries": [
    {
      "account_id": "<your-account-id>",
      "analysis_id": "<some-analysis-id>",
      "created_at": "2016-07-20T16:43:50.567Z",
      "description": "The project source is required to include a valid .about.yml file.",
      "duration": 6000,
      "id": "80c00035-3aff-cd42-a53d-cb42d64a5e95",
      "name": "Has a valid .about.yml file",
      "passed": true,
      "project_id": "<some-project-id>",
      "results": {
        "about_yml": {
          "content": "---\nname: sonar-auth-geoaxis\nfull_name: SonarQube GEOAxIS OAuth2 Java plugin\ndescription: |\n             OAuth plugin specialized for the GEOAxIS Oracle Access Manager's provider.  \n             The current version (1.0.0) is compatible with SonarQube 5.4 and will need to be updated prior to upgrading to SonarQube 5.5+.  \n             SonarQube 5.4 does not support authentication via OAuth and LDAP plugins as the same time, therefore this plugin will not function properly if installed with the LDAP plugin.\nimpact: Allows integration of SonarQube into GEOAxIS SSO solution \nowner_type: project\nstage: live\ntestable: true\nteam:\n- github: howellsd \n  role: lead \n  id: howellsd \nlicenses:\n  sonar-auth-geoaxis:\n    name: CC0 \n    url: https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis/blob/master/LICENSE.md \n",
          "message": "",
          "valid": true
        },
      "risk": "low",
      "summary": "Finished about_yml scan for Sonar-Plugin, valid .about.yml found.",
      "type": "about_yml",
      "updated_at": "2016-07-20T16:43:50.567Z"
    }
  ],
  "source": "https://gitlab.devops.geointservices.io/DevOps/sonar-auth-geoaxis.git",
  "status": "finished",
  "summary": "",
  "text": null,
  "trigger": "source commit",
  "trigger_author": "floydpepper",
  "trigger_hash": "d4e3gc6",
  "trigger_text": "fix breaking changes",
  "type": "git",
  "updated_at": "2016-07-20T16:43:58.309Z"
}
```
