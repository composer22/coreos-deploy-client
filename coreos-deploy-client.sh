#!/bin/bash

# Set base environment variables.
function init() {
	if [ -z "$COREOS_DEPLOY_URL" ]
	then
		echo "Enter the URL/port to the deploy service (Ex: http://foo.com:80) and press [ENTER]:"
	    read url
	    export COREOS_DEPLOY_URL=$url
	fi
	if [ -z "$COREOS_DEPLOY_URL" ]
	then
		printf "\e[0;31mDeploy URL is mandatory.\e[0m\n"
		exit 1
	fi

	if [ -z "$COREOS_DEPLOY_TOKEN" ]
	then
		echo "Enter your CoreOS Deploy API token and press [ENTER]:"
		read token
		export COREOS_DEPLOY_TOKEN=$token
	fi
	if [ -z "$COREOS_DEPLOY_TOKEN" ]
	then
		printf "\e[0;31mDeploy API token is mandatory.\e[0m\n"
		exit 1
	fi
}

# Deploy a service.
function deploy() {
	init
	name=
	version=
	instances=
	template_filepath=
	etcd2_filepath=
	shift
	while getopts "n:r:i::t:e::" VALUE "$@" ; do
	  	case "$VALUE" in
			n) name="$OPTARG" ;;
			r) version="$OPTARG" ;;
			i) instances="$OPTARG" ;;
			t) template_filepath="$OPTARG" ;;
			e) etcd2_filepath="$OPTARG" ;;
		esac
	done

	if [ -z "$name" ]
	then
		echo "Enter application name [ENTER]:"
		read inp
		name=$inp
	fi
	if [ -z "$name" ]
	then
		printf "\e[0;31mName is mandatory.\e[0m\n"
		skip_header=1
		namespace=help
		targeted_help=deploy
		return 0
	fi

	if [ -z "$version" ]
	then
		echo "Enter application version [ENTER]:"
		read inp
		version=$inp
	fi
	if [ -z "$version" ]
	then
		printf "\e[0;31mVersion is mandatory.\e[0m\n"
		skip_header=1
		namespace=help
		targeted_help=deploy
		return 0
	fi

	if [ -z "$instances" ]
	then
		echo "Enter optional instances (default = 2) [ENTER]:"
		read inp
		instances=$inp
	fi
	if [ -z "$instances" ]
	then
		instances=2
	fi

	if [ -z "$template_filepath" ]
	then
		echo "Enter .service template filepath [ENTER]:"
		read inp
		template_filepath=$inp
	fi
	if [ -z "$template_filepath" ]
	then
		printf "\e[0;31mService template is mandatory.\e[0m\n"
		skip_header=1
		namespace=help
		targeted_help=deploy
		return 0
	fi

	if [ -z "$etcd2_filepath" ]
	then
		echo "Enter optional etc2 key/value filepath [ENTER]:"
		read inp
		etcd2_filepath=$inp
	fi
	opt_etcd2=
	if [ ! -z "$etcd2_filepath" ]
	then
		opt_etcd2="-e $etcd2_filepath"
	fi

	./coreos-deploy-client -n $name -r $version -i $instances \
	  	-t $template_filepath $opt_etcd2 -u $COREOS_DEPLOY_URL \
		-b $COREOS_DEPLOY_TOKEN
	break
}

# Get status of a previous deploy.
function status() {
	init
	deploy_id=
	shift
	while getopts "p:" VALUE "$@" ; do
	  	case "$VALUE" in
			p) deploy_id="$OPTARG" ;;
		esac
	done

	if [ -z "$deploy_id" ]
	then
		echo "Enter deploy id to query [ENTER]:"
		read inp
		deploy_id=$inp
	fi
	if [ -z "$deploy_id" ]
	then
		printf "\e[0;31mDeploy ID is mandatory.\e[0m\n"
		skip_header=1
		namespace=help
		targeted_help=status
		return 0
	fi

	./coreos-deploy-client -u $COREOS_DEPLOY_URL -b $COREOS_DEPLOY_TOKEN \
		-p $deploy_id
	break
}

# General command line help.
function help_usage() {
  cat << EOF
$HELP_HEADER

Usage:
  coreos-deploy-client.sh (deploy|status|help) [options]

  Take help to list the command line options for each command.

  examples:

  coreos-deploy-client.sh help
  coreos-deploy-client.sh help deploy

  You might want to add the server URL and API token to your environment.
  Both are mandatory so you will be prompted if they are not found:

  export COREOS_DEPLOY_URL=<the server URL ex: http://foo.com:80>
  export COREOS_DEPLOY_TOKEN=<your API token>

EOF
}

# Menu options for help.
function help_index() {
	cat << EOF
$HELP_HEADER

Available help:
1) deploy
2) status
3) usage
EOF
echo "Pick what you need by typing a number: "
}

# Help for deploy command.
function help_deploy() {
  [ -z "$1" ] && echo $HELP_HEADER
  cat << EOF

Deploy is used to submit a deployment request to the cluster. It returns
a unique deploy_id that can be used as a reference for checking the status
of the deploy job.

Deploy command usage:

  coreos-deploy-client.sh deploy [options]

  -n    The name of the application service (ex mobile-video).

  -r    The unique version of the service (ex: 1.0.2).

  -i    The number of service instances to launch. Defaults to 2 instances.

  -t    The path to the .service file source code. This is the fleet unit
        file to launch the application in the cluster.
        ex: /path/to/mobile-video.service

  -e    The path to the etcd2 key/value file that contains etcd2 keys for
        update on the cluster. This is optional. The file should contains
		a key and a value on each line, seperated by a space. For example:

		key1 value1
		key2 value2
		keyn a string value

		Path example: /path/to/mobile-video.etc2keys

EOF
}

# Help for status command.
function help_status() {
  [ -z "$1" ] && echo $HELP_HEADER
  cat << EOF

Status is used to return the results of a previous deploy request.

Status command usage:

  coreos-deploy-client.sh status [options]

  -p    The deploy_id to query on the server.
        ex: 35E54E6B-51F4-455D-B431-91030E46735B

EOF
}

# Handle a help request.
function parse_help() {
	[[ -z "$targeted_help" ]] && targeted_help=$2
	if [ -z "$targeted_help" ]
	then
		skip_header=1
		help_index
	  	read targeted_help
	else
	  	case "$targeted_help" in
	    	1|deploy) help_deploy $skip_header ;;
	    	2|status) help_status $skip_header ;;
			3|usage) help_usage $skip_header ;;
	    	*) help_usage ;;
	  	esac
	  	break
	fi
}

################## MAINLINE ############
HELP_HEADER="CoreOS Deploy Client Command Line Interface"
namespace=$1
skip_header=
targeted_help=

while true
do
  case "$namespace" in
    deploy) deploy $@ ;;
    status) status $@ ;;
    help) parse_help $@ ;;
    *) help_usage && exit 1 ;;
  esac
done
