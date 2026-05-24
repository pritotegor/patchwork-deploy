// Package inventory provides types and helpers for loading and filtering
// the set of remote hosts that patchwork-deploy targets.
//
// An inventory file is a JSON document with the following shape:
//
//	{
//	  "hosts": [
//	    {
//	      "name":    "web-01",
//	      "address": "10.0.0.1",
//	      "user":    "deploy",
//	      "port":    22,
//	      "tags":    ["web", "prod"]
//	    }
//	  ]
//	}
//
// The user field defaults to "root" and port defaults to 22 when omitted.
package inventory
