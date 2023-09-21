package service

import (
	"context"

	"github.com/superjcd/rivalservice/config"
	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"github.com/superjcd/rivalservice/pkg/database"
	"github.com/superjcd/rivalservice/service/sql_store"
	"github.com/superjcd/rivalservice/service/sql_store/sql"

	"gorm.io/gorm"
)

var _DB *gorm.DB

// Server Server struct
type Server struct {
	v1.UnimplementedRivalServiceServer
	datastore sql_store.SqlFactory
	client    v1.RivalServiceClient
	conf      *config.Config
}

// NewServer New service grpc server
func NewServer(conf *config.Config, client v1.RivalServiceClient) (v1.RivalServiceServer, error) {
	_DB = database.MustPreParePostgresqlDb(&conf.Pg)
	factory, err := sql.NewSqlStoreFactory(_DB)
	if err != nil {
		return nil, err
	}

	server := &Server{
		client:    client,
		datastore: factory,
		conf:      conf,
	}

	return server, nil
}

func (s *Server) CreateRival(ctx context.Context, rq *v1.CreateRivalRequest) (*v1.CreateRivalResponse, error) {
	if err := s.datastore.Rivals().Create(ctx, rq); err != nil {
		return &v1.CreateRivalResponse{Msg: "创建失败", Status: v1.Status_failure}, err
	}
	return &v1.CreateRivalResponse{
		Msg:    "创建成功",
		Status: v1.Status_success,
	}, nil
}

func (s *Server) ListRival(ctx context.Context, rq *v1.ListRivalRequest) (*v1.ListRivalResponse, error) {
	groups, err := s.datastore.Rivals().List(ctx, rq)
	if err != nil {
		return &v1.ListRivalResponse{Msg: "获取列表失败", Status: v1.Status_failure}, err
	}

	resp := groups.ConvertToListRivalResponse("成功获取列表", v1.Status_success)

	return &resp, nil
}

func (s *Server) DeleteRival(ctx context.Context, rq *v1.DeleteRivalRequest) (*v1.DeleteRivalResponse, error) {
	if err := s.datastore.Rivals().Delete(ctx, rq); err != nil {
		return &v1.DeleteRivalResponse{Msg: "删除失败", Status: v1.Status_failure}, err
	}

	return &v1.DeleteRivalResponse{Msg: "删除成功", Status: v1.Status_success}, nil
}

func (s *Server) AppendRivalChanges(ctx context.Context, rq *v1.AppendRivalChangesRequest) (*v1.AppendRivalChangesResponse, error) {
	if err := s.datastore.RivalChanges().Append(ctx, rq); err != nil {
		return &v1.AppendRivalChangesResponse{Msg: "新增竞品变化失败", Status: v1.Status_failure}, err
	}

	return &v1.AppendRivalChangesResponse{Msg: "删除成功", Status: v1.Status_success}, nil
}

func (s *Server) ListRivalChanges(ctx context.Context, rq *v1.ListRivalChangesRequest) (*v1.ListRivalChangesResponse, error) {
	result, err := s.datastore.RivalChanges().List(ctx, rq)

	if err != nil {
		return &v1.ListRivalChangesResponse{Msg: "获取竞品变化列表失败", Status: v1.Status_failure}, err
	}

	resp := result.ConvertToListRivalChangeResponse("获取竞品变化列表成功", v1.Status_success)
	return &resp, nil

}

func (s *Server) DeleteRivalChanges(ctx context.Context, rq *v1.DeleteRivalChangesRequest) (*v1.DeleteRivalChangesResponse, error) {
	if err := s.datastore.RivalChanges().Delete(ctx, rq); err != nil {
		return &v1.DeleteRivalChangesResponse{
			Msg:    "删除竞争对手变化数据失败",
			Status: v1.Status_failure,
		}, err
	}
	return &v1.DeleteRivalChangesResponse{
		Msg:    "删除竞争对手变化数据成功",
		Status: v1.Status_success,
	}, nil
}

func (s *Server) AppendRivalProductInactiveDetail(ctx context.Context, rq *v1.AppendRivalProductInactiveDetailRequest) (*v1.AppendRivalProductInactiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendInactiveDetail(ctx, rq); err != nil {
		return &v1.AppendRivalProductInactiveDetailResponse{
			Msg:    "追加详情数据失败",
			Status: v1.Status_failure,
		}, err
	}

	return &v1.AppendRivalProductInactiveDetailResponse{
		Msg:    "追加详情数据成功",
		Status: v1.Status_failure,
	}, nil
}

func (s *Server) DeleteRivalInactiveDetail(ctx context.Context, rq *v1.DeleteRivalInactiveDetailRequest) (*v1.DeleteRivalInactiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().DeleteInactiveDetail(ctx, rq); err != nil {
		return &v1.DeleteRivalInactiveDetailResponse{Msg: "删除inactive产品详情数据失败", Status: v1.Status_failure}, err
	}
	return &v1.DeleteRivalInactiveDetailResponse{Msg: "删除inactive产品详情数据成功", Status: v1.Status_success}, nil
}

func (s *Server) AppendRivalProductActiveDetail(ctx context.Context, rq *v1.AppendRivalProductActiveDetailRequest) (*v1.AppendRivalProductActiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendActiveDetail(ctx, rq); err != nil {
		return &v1.AppendRivalProductActiveDetailResponse{
			Msg:    "追加详情数据失败",
			Status: v1.Status_failure,
		}, err
	}

	return &v1.AppendRivalProductActiveDetailResponse{
		Msg:    "追加详情数据成功",
		Status: v1.Status_failure,
	}, nil
}

func (s *Server) DeleteRivalActiveDetail(ctx context.Context, rq *v1.DeleteRivalActiveDetailRequest) (*v1.DeleteRivalActiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().DeleteActiveDetail(ctx, rq); err != nil {
		return &v1.DeleteRivalActiveDetailResponse{Msg: "删除inactive产品详情数据失败", Status: v1.Status_failure}, err
	}
	return &v1.DeleteRivalActiveDetailResponse{Msg: "删除inactive产品详情数据成功", Status: v1.Status_success}, nil
}
