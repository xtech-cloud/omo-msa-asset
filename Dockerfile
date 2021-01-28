FROM alpine:3.11
ADD omo.msa.asset /usr/bin/omo.msa.asset
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.asset" ]
