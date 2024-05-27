package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
	"strconv"
)

type FolderService struct{}

func switchFolder(info *cache.FolderInfo) *pb.FolderInfo {
	tmp := new(pb.FolderInfo)
	tmp.Uid = info.UID
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Owner = info.Scene
	tmp.Parent = info.Parent
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Type = uint32(info.Type)
	tmp.Access = uint32(info.Access)
	tmp.Count = info.GetChildCount()

	tmp.Tags = info.Tags
	tmp.Users = info.Users
	tmp.Contents = make([]*pb.PairInfo, 0, len(info.Contents))
	for _, content := range info.Contents {
		tmp.Contents = append(tmp.Contents, &pb.PairInfo{Key: content.Key, Value: content.Value, Count: content.Count})
	}
	return tmp
}

func (mine *FolderService) AddOne(ctx context.Context, in *pb.ReqFolderAdd, out *pb.ReplyFolderInfo) error {
	path := "folder.addOne"
	inLog(path, in)

	info, err := cache.Context().CreateFolder(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFolder(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFolderInfo) error {
	path := "folder.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	folder, err := cache.Context().GetFolder(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFolder(folder)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "folder.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the folder uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	err := cache.Context().RemoveFolder(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) UpdateBase(ctx context.Context, in *pb.ReqFolderUpdate, out *pb.ReplyInfo) error {
	path := "folder.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	folder, err := cache.Context().GetFolder(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	err = folder.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "folder.getStatistic"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) UpdateContents(ctx context.Context, in *pb.ReqFolderContents, out *pb.ReplyFolderInfo) error {
	path := "folder.updateContents"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	folder, err := cache.Context().GetFolder(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	err = folder.UpdateContents(in.Operator, in.Contents)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFolder(folder)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "folder.updateByFilter"
	inLog(path, in)
	folder, err := cache.Context().GetFolder(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	if in.Field == "parent" {
		err = folder.UpdateParent(in.Operator, in.Value)
	} else if in.Field == "cover" {
		err = folder.UpdateCover(in.Operator, in.Value)
	} else if in.Field == "append" {
		err = folder.AppendContent(in.Value, "", in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FolderService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyFolderList) error {
	path := "folder.getByFilter"
	inLog(path, in)
	var err error
	var list []*cache.FolderInfo
	if in.Key == "scene" || in.Key == "" {
		tp, _ := strconv.Atoi(in.Value)
		list, err = cache.Context().GetFoldersByScene(in.Owner, uint8(tp))
	} else if in.Key == "parent" {
		list, err = cache.Context().GetFoldersByParent(in.Value)
	} else if in.Key == "scenes" {
		tp, _ := strconv.Atoi(in.Value)
		list, err = cache.Context().GetFoldersByScenes(in.Values, uint8(tp))
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.FolderInfo, 0, len(list))
	for _, info := range list {
		tmp := switchFolder(info)
		out.List = append(out.List, tmp)
	}
	out.Status = outLog(path, out)
	return nil
}
