package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
)

type ThumbService struct {}

func switchThumb(info *cache.ThumbInfo) *pb.ThumbInfo {
	tmp := new(pb.ThumbInfo)
	tmp.Uid = info.UID
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Owner = info.Owner
	tmp.Asset = info.Asset
	tmp.Probably = info.Probably
	tmp.Face = info.Face
	tmp.Url = info.URL
	return tmp
}

func (mine *ThumbService)AddOne(ctx context.Context, in *pb.ReqThumbAdd, out *pb.ReplyThumbOne) error {
	path := "thumb.addOne"
	if len(in.Asset) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}

	asset := cache.Context().GetAsset(in.Asset)
	if asset == nil {
		out.Status = outError(path, "not found the asset", pb.ResultStatus_NotExisted)
		return nil
	}

	info,err := asset.CreateThumb(in.Face, in.Url, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchThumb(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyThumbOne) error {
	path := "thumb.getOne"
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	var thumb *cache.ThumbInfo
	if len(in.Owner) > 1 {
		asset := cache.Context().GetAsset(in.Owner)
		if asset == nil {
			out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
			return nil
		}
		thumb = asset.GetThumb(in.Uid)
	}else{
		thumb = cache.Context().GetThumb(in.Uid)
	}
	if thumb == nil {
		out.Status = outError(path, "the thumb not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchThumb(thumb)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "thumb.removeOne"
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the thumb uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the asset uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	asset := cache.Context().GetAsset(in.Owner)
	if asset == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := asset.RemoveThumb(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService)GetList(ctx context.Context, in *pb.ReqThumbList, out *pb.ReplyThumbList) error {
	path := "thumb.getList"
	if len(in.Asset) < 1 {
		out.Status = outError(path, "the asset uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	asset := cache.Context().GetAsset(in.Asset)
	if asset == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.ThumbInfo, 0, len(in.List))
	for _, val := range asset.Thumbs {
		out.List = append(out.List, switchThumb(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ThumbService)GetByOwner(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyThumbList) error {
	path := "thumb.getByOwner"
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	list := cache.Context().GetThumbsByOwner(in.Owner)
	out.Owner = in.Owner
	out.List = make([]*pb.ThumbInfo, 0, len(list))
	for _, val := range list {
		out.List = append(out.List, switchThumb(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ThumbService)UpdateBase(ctx context.Context, in *pb.ReqThumbUpdate, out *pb.ReplyThumbOne) error {
	path := "thumb.updateBase"
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	var thumb *cache.ThumbInfo
	if len(in.Owner) > 1 {
		asset := cache.Context().GetAsset(in.Owner)
		if asset == nil {
			out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
			return nil
		}
		thumb = asset.GetThumb(in.Uid)
	}else{
		thumb = cache.Context().GetThumb(in.Uid)
	}
	if thumb == nil {
		out.Status = outError(path, "the thumb not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := thumb.UpdateBase(in.Owner, in.Probably)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchThumb(thumb)
	out.Status = outLog(path, out)
	return nil
}

