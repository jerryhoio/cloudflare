package handler

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	proto "github.com/jerryhoio/cloudflare/pkg/proto"
)

type domainsHandler struct {
	cfCli *cloudflare.API
}

func NewDomainsHandler(cfClient *cloudflare.API) *domainsHandler {
	return &domainsHandler{cfCli: cfClient}
}

func (h *domainsHandler) List(ctx context.Context, req *proto.ListDomainsRequest, res *proto.ListDomainsResponse) error {
	zones, err := h.cfCli.ListZones()
	if err != nil {
		return err
	}
	var domains []*proto.Domain
	for _, z := range zones {
		var nameservers []string
		for _, n := range z.NameServers {
			nameservers = append(nameservers, n)
		}
		domains = append(domains, &proto.Domain{
			Name:        z.Name,
			Id:          z.ID,
			NameServers: nameservers,
			Status:      z.Status,
		})
	}
	res.Domains = domains
	return nil
}

func (h *domainsHandler) GetRecords(ctx context.Context, req *proto.GetDomainRecordsRequest, res *proto.GetDomainRecordsResponse) error {
	zoneId, err := h.cfCli.ZoneIDByName(req.DomainName)
	if err != nil {
		return err
	}
	filter := cloudflare.DNSRecord{}
	dnRec, err := h.cfCli.DNSRecords(zoneId, filter)

	if err != nil {
		return err
	}
	var records []*proto.DomainRecord
	for _, r := range dnRec {
		records = append(records, &proto.DomainRecord{
			Name:    r.Name,
			Type:    r.Type,
			Content: r.Content,
			Ttl:     int32(r.TTL),
		})
	}
	res.Records = records
	return nil
}
