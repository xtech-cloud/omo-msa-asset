package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-asset/proto/asset"
	"omo.msa.asset/cache"
	"strings"
)

type LabelService struct{}

func switchLabel(info *cache.LabelInfo) *pb.LabelInfo {
	tmp := new(pb.LabelInfo)
	tmp.Uid = info.UID
	tmp.Updated = info.Updated
	tmp.Created = info.Created
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator

	tmp.Scene = info.Scene
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Type = uint32(info.Type)
	return tmp
}

func (mine *LabelService) AddOne(ctx context.Context, in *pb.ReqLabelAdd, out *pb.ReplyLabelInfo) error {
	path := "label.addOne"
	inLog(path, in)

	in.Name = strings.TrimSpace(in.Name)
	had, er := cache.Context().HadLabel(in.Name)
	if er != nil {
		out.Status = outError(path, er.Error(), pb.ResultStatus_DBException)
		return nil
	}
	if had {
		out.Status = outError(path, "the name had existed", pb.ResultStatus_Repeated)
		return nil
	}

	info, err := cache.Context().CreateLabel(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchLabel(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LabelService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyLabelInfo) error {
	path := "label.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the asset is empty", pb.ResultStatus_Empty)
		return nil
	}
	label, err := cache.Context().GetLabel(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchLabel(label)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LabelService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "label.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the label uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	err := cache.Context().RemoveLabel(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, "the asset not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *LabelService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "label.getStatistic"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty", pb.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *LabelService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "label.updateByFilter"
	inLog(path, in)
	info, err := cache.Context().GetLabel(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	if in.Field == "name" {
		err = info.UpdateBase(in.Value, info.Remark, in.Operator)
	} else if in.Field == "remark" {
		err = info.UpdateBase(info.Name, in.Value, in.Operator)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *LabelService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyLabelList) error {
	path := "label.getByFilter"
	inLog(path, in)
	var err error
	var list []*cache.LabelInfo
	if in.Key == "scene" || in.Key == "" {
		list, err = cache.Context().GetLabelsByScene(in.Owner)
	} else if in.Key == "asset_quote" {
		list, err = cache.Context().GetLabelsByQuote(in.Value)
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.LabelInfo, 0, len(list))
	for _, info := range list {
		tmp := switchLabel(info)
		out.List = append(out.List, tmp)
	}
	out.Status = outLog(path, out)
	return nil
}
