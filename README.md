# omo-msa-asset
微服务-资源

MICRO_REGISTRY=consul micro call omo.msa.asset AssetService.AddOne '{"name":"tese1", "md5":"11111", "owner":"hzz", "type":1, "size":500, "language":"zh", "version":"222222"}'
MICRO_REGISTRY=consul micro call omo.msa.asset AssetService.GetByOwner '{"owner":"hzz"}'
