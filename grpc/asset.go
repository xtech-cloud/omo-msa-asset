package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
	"omo.msa.asset/config"
	"strconv"
)

type AssetService struct{}

func switchAsset(owner string, info *cache.AssetInfo) *pb.AssetInfo {
	tmp := new(pb.AssetInfo)
	tmp.Owner = owner
	tmp.Uid = info.UID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Meta = info.Meta
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	if info.Owner != owner {
		tmp.Owner = info.Owner
	}

	tmp.Language = info.Language
	tmp.Type = uint32(info.Type)
	tmp.Scope = uint32(info.Scope)
	tmp.Size = info.Size
	tmp.Url = info.URL()
	tmp.Format = info.Format
	tmp.Md5 = info.MD5
	tmp.Version = info.Version
	tmp.Snapshot = info.SnapshotURL()
	tmp.Small = info.SmallImageURL()
	tmp.Width = info.Width
	tmp.Height = info.Height
	tmp.Weight = info.Weight
	tmp.Quote = info.Quote
	tmp.Status = uint32(info.Status)
	tmp.Links = info.Links
	tmp.Tags = info.Tags
	tmp.Source = info.SourceURL()

	thumbs, er := info.GetThumbs()
	if er == nil {
		tmp.Thumbs = make([]*pb.ThumbBrief, 0, len(thumbs))
		for _, thumb := range thumbs {
			tmp.Thumbs = append(tmp.Thumbs, switchThumbBrief(thumb))
		}
	}

	return tmp
}

func switchThumbBrief(info *cache.ThumbInfo) *pb.ThumbBrief {
	tmp := new(pb.ThumbBrief)
	tmp.Uid = info.UID
	tmp.Owner = info.Owner
	tmp.Face = info.User
	tmp.Blur = info.Blur
	tmp.Url = cache.GetURL(info.File, true)
	tmp.Similar = info.Similar
	tmp.Probably = info.Probably
	return tmp
}

func (mine *AssetService) AddOne(ctx context.Context, in *pb.ReqAssetAdd, out *pb.ReplyAssetOne) error {
	path := "asset.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pb.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateAsset(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAsset(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetOne) error {
	path := "asset.getOne"
	//inLog(path, in)
	var info *cache.AssetInfo
	if len(in.Uid) < 1 {
		info = cache.Context().GetAssetByKey(in.Operator)
	} else {
		info = cache.Context().GetAsset(in.Uid)
	}
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchAsset(in.Owner, info)
	out.Status = outNonLog()
	return nil
}

func (mine *AssetService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "asset.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	var err error
	info := cache.Context().GetAsset(in.Uid)
	if info != nil {
		err = info.Remove(in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) GetList(ctx context.Context, in *pb.ReqAssetList, out *pb.ReplyAssetList) error {
	path := "asset.getList"
	inLog(path, in)
	out.List = make([]*pb.AssetInfo, 0, len(in.List))
	for _, val := range in.List {
		info := cache.Context().GetAsset(val)
		if info != nil {
			out.List = append(out.List, switchAsset(info.Owner, info))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AssetService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "asset.getStatistic"
	inLog(path, in)
	if in.Key == "quote" {
		out.Count = cache.Context().GetAssetCount(in.Value)
	} else if in.Key == "quote_creator" {
		out.Count = cache.Context().GetAssetCountByQuoteCreator(in.Value, in.Operator)
	} else if in.Key == "owner_creator" {
		out.Count = cache.Context().GetAssetCountByOwnerCreator(in.Value, in.Operator)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AssetService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAssetList) error {
	path := "asset.getByFilter"
	inLog(path, in)
	var list []*cache.AssetInfo
	var total uint32
	var pages uint32
	if in.Key == "links" {
		list = cache.Context().GetAssetsByLink(in.Key)
	} else if in.Key == "creator" {
		list = cache.Context().GetAssetsByCreator(in.Value)
	} else if in.Key == "type" {
		tp, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetAssetsByType(tp)
	} else if in.Key == "regex" {
		if len(in.Numbers) == 2 {
			list = cache.Context().GetAssetsByRegex(in.Value, in.Numbers[0], in.Numbers[1])
		}
	} else if in.Key == "quote" {
		total, pages, list = cache.Context().GetAssetsByQuote(in.Owner, in.Value, in.Page, in.Number)
	} else if in.Key == "quote_creator" {
		total, pages, list = cache.Context().GetAssetsByQuoteCreator(in.Value, in.Operator, in.Page, in.Number)
	} else if in.Key == "owner_type" {
		tp, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetAssetsByOwnerType(in.Owner, tp)
	} else if in.Key == "owners_type" {
		tp, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
			return nil
		}
		list = cache.Context().GetAssetsByOwnersType(in.Values, tp)
	} else if in.Key == "quote_status" {
		st, err := strconv.Atoi(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
			return nil
		}
		total, pages, list = cache.Context().GetAssetsByQuoteStatus(in.Owner, uint32(st), in.Page, in.Number)
	}
	out.List = make([]*pb.AssetInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchAsset(info.Owner, info))
	}
	out.Total = total
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AssetService) GetByOwner(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetList) error {
	path := "asset.getByOwner"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	var list []*cache.AssetInfo
	if in.Operator == "publish" {
		var st = uint32(cache.StatusPublish)
		if len(in.Uid) > 0 {
			ss, _ := strconv.ParseUint(in.Uid, 10, 32)
			if ss == uint64(cache.StatusVisible) {
				st = uint32(ss)
			}
		}
		list = cache.Context().GetPublishAssetsByOwner(in.Owner, st)
	} else {
		list = cache.Context().GetAssetsByOwner(in.Owner)
	}

	out.Owner = in.Owner
	out.List = make([]*pb.AssetInfo, 0, len(list))
	for _, val := range list {
		out.List = append(out.List, switchAsset(in.Owner, val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AssetService) GetToken(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetToken) error {
	path := "asset.getToken"
	inLog(path, in)
	out.Expire = uint32(config.Schema.Storage.Expire)
	out.Domain = config.Schema.Storage.Domain
	out.Bucket = config.Schema.Storage.Bucket
	out.Limit = uint32(config.Schema.Storage.Limit)
	out.Access = config.Schema.Storage.AccessKey
	out.Secret = config.Schema.Storage.SecretKey
	out.Token = cache.Context().GetUpToken(in.Uid)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateSnapshot(ctx context.Context, in *pb.ReqAssetFlag, out *pb.ReplyInfo) error {
	path := "asset.updateSnapshot"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateSnapshot(in.Operator, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateSmall(ctx context.Context, in *pb.ReqAssetFlag, out *pb.ReplyInfo) error {
	path := "asset.updateSmall"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateSmall(in.Operator, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateBase(ctx context.Context, in *pb.ReqAssetUpdate, out *pb.ReplyInfo) error {
	path := "asset.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Operator, in.Name, in.Remark)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateMeta(ctx context.Context, in *pb.ReqAssetFlag, out *pb.ReplyInfo) error {
	path := "asset.updateMeta"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMeta(in.Operator, in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateWeight(ctx context.Context, in *pb.ReqAssetWeight, out *pb.ReplyInfo) error {
	path := "asset.updateWeight"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateWeight(in.Weight, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateStatus(ctx context.Context, in *pb.ReqAssetWeight, out *pb.ReplyInfo) error {
	path := "asset.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetAsset(in.Uid)
	if info == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(uint8(in.Weight), in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "asset.updateByFilter"
	inLog(path, in)
	var err error
	if in.Uid == "" {
		if in.Field == "status" {
			st, er := strconv.Atoi(in.Value)
			if er != nil {
				out.Status = outError(path, er.Error(), pb.ResultStatus_DBException)
			}
			err = cache.Context().UpdateAssetsStatus(in.Values, uint32(st), in.Operator)
		} else if in.Field == "publish" {
			err = cache.Context().PublishAssetsEntity(in.Value, in.Operator)
		} else if in.Field == "batch_scope" {
			err = cache.Context().BatchUpdateScope(in.Values)
		} else {
			out.Status = outError(path, "not define the field", pb.ResultStatus_DBException)
			return nil
		}
	} else {
		info := cache.Context().GetAsset(in.Uid)
		if info == nil {
			out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
			return nil
		}

		if in.Field == "type" {
			tp, _ := strconv.ParseUint(in.Value, 10, 32)
			err = info.UpdateType(uint8(tp), in.Operator)
		} else if in.Field == "language" {
			err = info.UpdateLanguage(in.Value, in.Operator)
		} else if in.Field == "links" {
			err = info.UpdateLinks(in.Operator, in.Values)
		} else if in.Field == "owner" {
			err = info.UpdateOwner(in.Operator, in.Value)
		} else if in.Field == "quote" {
			err = info.UpdateQuote(in.Operator, in.Value)
		} else if in.Field == "tags" {
			err = info.UpdateTags(in.Operator, in.Values)
		} else {
			out.Status = outError(path, "not define the field", pb.ResultStatus_DBException)
			return nil
		}
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}
