#!/bin/bash

ACTION_COUNT="count"
ACTION_CLEAN="clean"

ACTION=$1
DELETE_NUMBER=10

while getopts ":n:h" optname
do
    case "$optname" in
      "n")
        DELETE_NUMBER=$OPTARG
        ;;
      "h")
        printf "Usage:\n# count the si2runs with job without pod \n./s2iruns.sh count \n# clean the si2runs with job without pod \n./s2iruns.sh clean -n 100 \n    -n: the number of cleaning s2iruns\n"
        exit 0
        ;;
      ":")
        echo "No argument value for option $OPTARG"
        ;;
      "?")
        ;;
      *)
        echo "Unknown error while processing options"
        ;;
    esac
done

if [[ "${ACTION}" != "${ACTION_COUNT}" ]] && [[ "${ACTION}" != "${ACTION_CLEAN}" ]]; then
  echo "un-support action $ACTION!!"
  exit 1
fi




POD_NAME=""
operate() {
  NAMESPACE=$1
  RUN=$2
  STATE=$3
  JOB=$4

  if [[ "${STATE}" != "Successful" ]] && [[ "${STATE}" != "Failed" ]]; then
    echo "the s2irun(${RUN}) not complete, ignore .. "
    return 0
  fi

  if [ ! -n "${JOB}" ]; then
    echo "there is no job for s2irun(${RUN}), ignore .. "
    return 0
  fi

  getpod ${NAMESPACE} ${JOB}
  if [ ! -n "${POD_NAME}" ]; then
    COUNT=$(expr $COUNT + 1)
    if [[ "${ACTION}" == "${ACTION_CLEAN}" ]]; then
      echo "!!there is no pod for s2irun(${RUN}/${JOB})"
      echo "[`date "+%Y-%m-%d %H:%M:%S"`] deleting s2irun ${RUN}/${JOB}" >> ~/s2irun-clean.log
      kubectl -n ${NAMESPACE} delete s2iruns.devops.kubesphere.io ${RUN}
    fi
  fi
}


getpod() {
  ns=$1
  job=$2
  POD_NAME=""
  while read pod; do
    pod_job=$(echo ${pod} | awk '{print $7}')
    pod_ns=$(echo ${pod} | awk '{print $1}')
    if [[ "$job" == "$pod_job" ]] && [[ "$ns" == "$pod_ns" ]]; then
      POD_NAME=$(echo ${pod} | awk '{print $2}')
      break
    fi
  done <$POD_FILE
}




S2IRUN_FILE=/tmp/s2iruns.txt
POD_FILE=/tmp/pods.txt
kubectl get s2iruns.devops.kubesphere.io -A --no-headers=true --ignore-not-found=true > ${S2IRUN_FILE}
kubectl get pod -A --label-columns=job-name --no-headers=true --ignore-not-found=true | awk '{if (length($7) > 0) print $0}' > ${POD_FILE}

COUNT=0
POD_NAME=""
while read s2irun; do
  operate ${s2irun}
  if [ $COUNT -ge $DELETE_NUMBER ]; then
    exit 0
  fi
done <$S2IRUN_FILE

if [[ "${ACTION}" == "${ACTION_CLEAN}" ]]; then
  echo "clean count: ${COUNT}"
else
  echo "count: ${COUNT}"
fi
