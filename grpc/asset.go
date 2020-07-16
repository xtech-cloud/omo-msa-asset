package grpc

import (
	"context"
	"errors"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
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
	tmp.Uuid = info.UUID
	tmp.Format = info.Format
	tmp.Md5 = info.MD5
	tmp.Version = info.Version
	tmp.Language = info.Language
	return tmp
}

func (mine *AssetService)AddOne(ctx context.Context, in *pb.ReqAssetAdd, out *pb.ReplyAssetOne) error {
	if len(in.Owner) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	owner := cache.GetOwner(in.Owner)
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
	err := owner.CreateAsset(info)
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}else{
		out.Info = switchAsset(in.Owner, info)
	}
	return err
}

func (mine *AssetService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetOne) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the asset is empty")
	}
	if len(in.Owner) < 1 {
		info := cache.GetAsset(in.Uid)
		if info == nil {
			out.ErrorCode = pb.ResultStatus_NotExisted
			return errors.New("the asset not found")
		}
		out.Info = switchAsset(in.Owner, info)
	}else{
		owner := cache.GetOwner(in.Owner)
		info := owner.GetAsset(in.Uid)
		if info == nil {
			out.ErrorCode = pb.ResultStatus_NotExisted
			return errors.New("the asset not found")
		}
		out.Info = switchAsset(in.Owner, info)
	}

	return nil
}

func (mine *AssetService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetOne) error {

	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the asset is empty")
	}
	var err error
	if len(in.Owner) < 1 {
		err = cache.RemoveAsset(in.Uid, in.Operator)
	}else{
		owner := cache.GetOwner(in.Owner)
		err = owner.RemoveAsset(in.Uid, in.Operator)
	}
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	return err
}

func (mine *AssetService)GetList(ctx context.Context, in *pb.ReqAssetList, out *pb.ReplyAssetList) error {
	out.List = make([]*pb.AssetInfo, 0, len(in.List))
	for _, val := range in.List {
		info := cache.GetAsset(val)
		if info != nil {
			out.List = append(out.List, switchAsset(info.Owner, info))
		}
	}
	return nil
}

func (mine *AssetService)GetByOwner(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAssetList) error {
	if len(in.Owner) < 1 {
		return errors.New("the owner is empty")
	}
	owner := cache.GetOwner(in.Owner)
	if owner == nil {
		return errors.New("the owner not found")
	}
	out.Owner = in.Owner
	out.List = make([]*pb.AssetInfo, 0, len(owner.Assets()))
	for _, val := range owner.Assets() {
		out.List = append(out.List, switchAsset(in.Owner, val))
	}
	return nil
}


