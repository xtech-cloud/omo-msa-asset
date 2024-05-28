package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
)

type ThumbService struct{}

func switchThumb(info *cache.ThumbInfo) *pb.ThumbInfo {
	tmp := new(pb.ThumbInfo)
	tmp.Uid = info.UID
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Owner = info.Owner
	tmp.Asset = info.Asset
	tmp.Probably = info.Probably
	tmp.Url = cache.GetURL(info.File, true)
	tmp.Blur = info.Blur
	tmp.Similar = info.Similar
	tmp.Meta = info.Meta
	tmp.User = info.User
	tmp.Quote = info.Quote
	return tmp
}

func (mine *ThumbService) AddOne(ctx context.Context, in *pb.ReqThumbAdd, out *pb.ReplyThumbOne) error {
	path := "thumb.addOne"
	inLog(path, in)
	if len(in.Asset) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}

	asset := cache.Context().GetAsset(in.Asset)
	if asset == nil {
		out.Status = outError(path, "not found the asset", pb.ResultStatus_NotExisted)
		return nil
	}
	//if asset.HadThumbByFace(in.Face) {
	//	out.Status = outError(path, "the face repeated", pb.ResultStatus_Repeated)
	//	return nil
	//}
	info, err := asset.CreateThumb(in.Url, in.Operator, in.Owner, in.Probably, in.Similar, in.Blur)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchThumb(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyThumbOne) error {
	path := "thumb.getOne"
	inLog(path, in)
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
	} else {
		thumb = cache.Context().GetThumb(in.Uid)
	}
	if thumb == nil {
		out.Status = outError(path, "the thumb not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchThumb(thumb)
	out.Status = outLog(path, out.Info.Uid)
	return nil
}

func (mine *ThumbService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "thumb.removeOne"
	inLog(path, in)
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

func (mine *ThumbService) GetList(ctx context.Context, in *pb.ReqThumbList, out *pb.ReplyThumbList) error {
	path := "thumb.getList"
	inLog(path, in)
	if len(in.Asset) < 1 {
		out.Status = outError(path, "the asset uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	asset := cache.Context().GetAsset(in.Asset)
	if asset == nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	array, err := asset.GetThumbs()
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ThumbInfo, 0, len(in.List))
	for _, val := range array {
		out.List = append(out.List, switchThumb(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ThumbService) GetByOwner(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyThumbList) error {
	path := "thumb.getByOwner"
	inLog(path, in)
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

func (mine *ThumbService) UpdateBase(ctx context.Context, in *pb.ReqThumbUpdate, out *pb.ReplyThumbOne) error {
	path := "thumb.updateBase"
	inLog(path, in)
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
	} else {
		thumb = cache.Context().GetThumb(in.Uid)
	}
	if thumb == nil {
		out.Status = outError(path, "the thumb not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := thumb.UpdateBase(in.Owner, in.Similar)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchThumb(thumb)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "thumb.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 2 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	thumb := cache.Context().GetThumb(in.Uid)
	var err error
	if in.Field == "meta" {
		err = thumb.UpdateInfo(in.Value, in.Operator)
	} else if in.Field == "bind" {
		err = thumb.BindEntity(in.Value, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ThumbService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyThumbList) error {
	path := "thumb.getByFilter"
	inLog(path, in)
	var list []*cache.ThumbInfo
	if in.Key == "quote_users" {
		list = cache.Context().GetUserThumbsByQuote(in.Value)
	}
	out.List = make([]*pb.ThumbInfo, 0, len(list))
	for _, item := range list {
		out.List = append(out.List, switchThumb(item))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}
