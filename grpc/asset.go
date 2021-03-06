package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
	"omo.msa.asset/config"
)

type AssetService struct {}

func switchAsset(owner string, info *cache.AssetInfo) *pb.AssetInfo {
	tmp := new(pb.AssetInfo)
	tmp.Owner = owner
	tmp.Uid = info.UID
	tmp.Name = info.Name
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Language = info.Language
	tmp.Owner = info.Owner
	tmp.Type = pb.OwnerType(info.Type)
	tmp.Size = info.Size
	tmp.Uuid = info.URL()
	tmp.Format = info.Format
	tmp.Md5 = info.MD5
	tmp.Version = info.Version
	tmp.Snapshot = info.SnapshotURL()
	tmp.Small = info.SmallImageURL()
	tmp.Width = info.Width
	tmp.Height = info.Height
	thumbs,er := info.GetThumbs()
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
	tmp.Face = info.Face
	tmp.Blur = info.Blur
	tmp.Url = info.URL
	tmp.Similar = info.Similar
	tmp.Probably = info.Probably
	return tmp
}

func (mine *AssetService)AddOne(ctx context.Context, in *pb.ReqAssetAdd, out *pb.ReplyAssetOne) error {
	path := "asset.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	owner := cache.Context().GetOwner(in.Owner)
	info := new(cache.AssetInfo)
	info.Name = in.Name
	info.MD5 = in.Md5
	info.Format = in.Format
	info.Version = in.Version
	info.Owner = in.Owner
	info.Type = uint8(in.Type)
	info.UUID = in.Uuid
	info.Size = in.Size
	info.Language = in.Language
	info.Creator = in.Operator
	info.Small = in.Small
	info.Snapshot = in.Snapshot
	info.Width = in.Width
	info.Height = in.Height
	err := owner.CreateAsset(info)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAsset(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetOne) error {
	path := "asset.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	if len(in.Owner) < 1 {
		info := cache.Context().GetAsset(in.Uid)
		if info == nil {
			out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
			return nil
		}
		out.Info = switchAsset(in.Owner, info)
	}else{
		owner := cache.Context().GetOwner(in.Owner)
		info := owner.GetAsset(in.Uid)
		if info == nil {
			out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
			return nil
		}
		out.Info = switchAsset(in.Owner, info)
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "asset.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	var err error
	if len(in.Owner) < 1 {
		err = cache.Context().RemoveAsset(in.Uid, in.Operator)
	}else{
		owner := cache.Context().GetOwner(in.Owner)
		err = owner.RemoveAsset(in.Uid, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService)GetList(ctx context.Context, in *pb.ReqAssetList, out *pb.ReplyAssetList) error {
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

func (mine *AssetService)GetByOwner(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetList) error {
	path := "asset.getByOwner"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	owner := cache.Context().GetOwner(in.Owner)
	if owner == nil {
		out.Status = outError(path, "the owner not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Owner = in.Owner
	out.List = make([]*pb.AssetInfo, 0, len(owner.Assets()))
	for _, val := range owner.Assets() {
		out.List = append(out.List, switchAsset(in.Owner, val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AssetService)GetToken(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetToken) error {
	path := "asset.getToken"
	inLog(path, in)
	out.Expire = uint32(config.Schema.Storage.Expire)
	out.Domain = config.Schema.Storage.Domain
	out.Bucket = config.Schema.Storage.Bucket
	out.Limit = uint32(config.Schema.Storage.Limit)
	out.Token = cache.Context().GetUpToken()
	out.Status = outLog(path, out)
	return nil
}

func (mine *AssetService)UpdateSnapshot(ctx context.Context, in *pb.ReqAssetFlag, out *pb.ReplyInfo) error {
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

func (mine *AssetService)UpdateSmall(ctx context.Context, in *pb.ReqAssetFlag, out *pb.ReplyInfo) error {
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




