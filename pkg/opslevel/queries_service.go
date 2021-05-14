package opslevel

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/opslevel/opslevel-go"
	"github.com/shurcooL/graphql"
	"strings"
)

// Queries

type Service struct {
	opslevel.Service
	Tags ServiceTagsConnection `json:"tags"`
}

type ServiceTagsConnection struct {
	Nodes []opslevel.Tag
}

type ServiceConnection struct {
	Nodes    []Service
	PageInfo opslevel.PageInfo
}

// GetByAlias

func GetServiceWithAlias(client *opslevel.Client, alias string) (*Service, error) {
	var q struct {
		Account struct {
			Service Service `graphql:"service(alias: $service)"`
		}
	}
	v := opslevel.PayloadVariables{
		"service": graphql.String(alias),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, err
	}

	return &q.Account.Service, nil
}

// GetByID

func GetServiceWithId(client *opslevel.Client, id string) (*Service, error) {
	var q struct {
		Account struct {
			Service Service `graphql:"service(id: $service)"`
		}
	}
	v := opslevel.PayloadVariables{
		"service": graphql.ID(id),
	}
	if err := client.Query(&q, v); err != nil {
		return nil, fmt.Errorf("could not GetServiceWithId: %v using query: %s", err, pretty.Sprint(q))
	}

	return &q.Account.Service, nil
}

// By Framework

type ListServicesByFrameworkQuery struct {
	Account struct {
		Services ServiceConnection `graphql:"services(framework: $framework, after: $after, first: $first)"`
	}
}

func (q *ListServicesByFrameworkQuery) Query(client *opslevel.Client, framework string) error {
	var subQ ListServicesByFrameworkQuery
	v := opslevel.PayloadVariables{
		"framework": graphql.String(framework),
		"after":     q.Account.Services.PageInfo.End,
		"first":     graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, framework)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByFramework(client *opslevel.Client, framework string) ([]Service, error) {
	q := ListServicesByFrameworkQuery{}
	if err := q.Query(client, framework); err != nil {
		return []Service{}, err
	}
	return q.Account.Services.Nodes, nil
}

// By Language

type ListServicesByLanguageQuery struct {
	Account struct {
		Services ServiceConnection `graphql:"services(language: $language, after: $after, first: $first)"`
	}
}

func (q *ListServicesByLanguageQuery) Query(client *opslevel.Client, language string) error {
	var subQ ListServicesByLanguageQuery
	v := opslevel.PayloadVariables{
		"language": graphql.String(language),
		"after":    q.Account.Services.PageInfo.End,
		"first":    graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, language)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByLanguage(client *opslevel.Client, framework string) ([]Service, error) {
	q := ListServicesByLanguageQuery{}
	if err := q.Query(client, framework); err != nil {
		return []Service{}, err
	}
	return q.Account.Services.Nodes, nil
}

// By OwnerAlias

type ListServicesByOwnerAliasQuery struct {
	Account struct {
		Services ServiceConnection `graphql:"services(ownerAlias: $ownerAlias, after: $after, first: $first)"`
	}
}

func (q *ListServicesByOwnerAliasQuery) Query(client *opslevel.Client, ownerAlias string) error {
	var subQ ListServicesByOwnerAliasQuery
	v := opslevel.PayloadVariables{
		"ownerAlias": graphql.String(ownerAlias),
		"after":      q.Account.Services.PageInfo.End,
		"first":      graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, ownerAlias)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func ListServicesByOwnerAlias(client *opslevel.Client, framework string) ([]Service, error) {
	q := ListServicesByOwnerAliasQuery{}
	if err := q.Query(client, framework); err != nil {
		return []Service{}, err
	}
	return q.Account.Services.Nodes, nil
}

// By Tag

type listServicesByTagQuery interface {
	Query(client *opslevel.Client, key, value string) (error)
	Services() []Service
}

func ListServicesByTag(client *opslevel.Client, value string) ([]Service, error) {
	tagKV := strings.Split(value, ":")
	if len(tagKV) != 2 {
		return nil, fmt.Errorf("tag filter requires `value` in format 'key:value', or `key:` if only filtering on the presence of a tag.")
	}

	var query  listServicesByTagQuery
	if tagKV[1] == "" {
		query = &ListServicesByTagQuery{}
	} else {
		query = &ListServicesByTagValueQuery{}
	}

	err := query.Query(client, tagKV[0], tagKV[1])
	if err != nil {
		return nil, err
	}
	return query.Services(), nil
}


type ListServicesByTagQuery struct {
	Account struct {
		Services ServiceConnection `graphql:"services(tag: {key:$key, value:}, after: $after, first: $first)"`
	}
}

func (q *ListServicesByTagQuery) Services() []Service {
	return q.Account.Services.Nodes
}

func (q *ListServicesByTagQuery) Query(client *opslevel.Client, key, value string) error {
	var subQ ListServicesByTagQuery
	v := opslevel.PayloadVariables{
		"key":   graphql.String(key),
		"after": q.Account.Services.PageInfo.End,
		"first": graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, key, value)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

type ListServicesByTagValueQuery struct {
	Account struct {
		Services ServiceConnection `graphql:"services(tag: {key:$key, value:$value}, after: $after, first: $first)"`
	}
}

func (q ListServicesByTagValueQuery) Query(client *opslevel.Client, key, value string) error {
	var subQ ListServicesByTagValueQuery
	v := opslevel.PayloadVariables{
		"key":   graphql.String(key),
		"value":   graphql.String(value),
		"after": q.Account.Services.PageInfo.End,
		"first": graphql.Int(100),
	}
	if err := client.Query(&subQ, v); err != nil {
		return err
	}
	if subQ.Account.Services.PageInfo.HasNextPage {
		subQ.Query(client, key, value)
	}
	for _, service := range subQ.Account.Services.Nodes {
		q.Account.Services.Nodes = append(q.Account.Services.Nodes, service)
	}
	return nil
}

func (q *ListServicesByTagValueQuery) Services() []Service {
	return q.Account.Services.Nodes
}

