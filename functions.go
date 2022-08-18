package autowire

import (
	"fmt"
	"github.com/gomelon/meta"
	"go/types"
	"golang.org/x/tools/go/packages"
	"sort"
	"text/template"
)

type functions struct {
	pkgParser  *meta.PkgParser
	metaParser *meta.Parser
	pkg        *packages.Package
	pkgPath    string
}

func NewFunctions(gen *meta.TmplPkgGen) *functions {
	return &functions{
		pkg:        gen.PkgParser().Package(gen.PkgPath()),
		pkgPath:    gen.PkgPath(),
		pkgParser:  gen.PkgParser(),
		metaParser: gen.MetaParser(),
	}
}

func (f *functions) FuncMap() template.FuncMap {
	return map[string]any{
		"parseWire": f.ParseWire,
	}
}

type ParsedWireResult struct {
	Providers []types.Object
	Bindings  []*ProviderHolder
}

func (r *ParsedWireResult) HasProvider() bool {
	return len(r.Providers) > 0
}

type ProviderHolder struct {
	Provider      types.Object
	ProviderType  types.Type
	InjectedIface types.Type
	Order         int32
	IsBase        bool
}

//ParseWire
//1. 获取所有ProviderFunc
//2. 获取ProviderType
//3. 获取ProviderTypeInterface
//4. 按ProviderTypeInterface分组
//5. 按order给各分组排序
//6. 模版输出
func (f *functions) ParseWire() (result *ParsedWireResult, err error) {
	result = &ParsedWireResult{}
	pkgFunctions := f.pkgParser.Functions(f.pkgPath)
	if len(pkgFunctions) == 0 {
		return
	}

	result.Providers = f.metaParser.FilterByMeta(MetaWireProvider, pkgFunctions)

	providerTypeIfaceToProviders := map[types.Type][]types.Object{}
	for _, function := range result.Providers {
		providerObj := f.pkgParser.FirstResult(function)
		if providerObj == nil {
			err = fmt.Errorf("provider function expect has one result but none,function=%s", function.String())
			return
		}

		providerType := providerObj.(*types.Var).Type()
		providerInterfaces := f.pkgParser.AnonymousAssignTo(providerType)
		if len(providerInterfaces) == 0 {
			continue
		}

		for _, iface := range providerInterfaces {
			providerTypeIfaceToProviders[iface] = append(providerTypeIfaceToProviders[iface], function)
		}

		firstParam := f.pkgParser.FirstParam(function)
		if firstParam == nil {
			continue
		}
	}

	for providerTypeIface, itfProviders := range providerTypeIfaceToProviders {
		itfProviderHolders := make([]*ProviderHolder, 0, len(itfProviders))
		for _, provider := range itfProviders {
			providerObject := f.pkgParser.FirstResult(provider)
			params := f.pkgParser.Params(provider)
			needInjectParam := false
			for _, param := range params {
				if !f.pkgParser.AssignableTo(param.Type(), providerTypeIface) {
					continue
				}
				needInjectParam = true
				wireMeta := f.metaParser.ObjectMetaGroup(provider, MetaWireProvider)[0]
				providerHolder := &ProviderHolder{
					Provider:      provider,
					ProviderType:  providerObject.Type(),
					InjectedIface: param.Type(),
					Order:         Order(wireMeta),
				}
				itfProviderHolders = append(itfProviderHolders, providerHolder)
			}
			if !needInjectParam {
				providerHolder := &ProviderHolder{
					Provider:      provider,
					ProviderType:  providerObject.Type(),
					InjectedIface: providerTypeIface,
					IsBase:        true,
				}
				itfProviderHolders = append(itfProviderHolders, providerHolder)
			}
		}
		sort.Slice(itfProviderHolders, func(i, j int) bool {
			return itfProviderHolders[i].IsBase ||
				!itfProviderHolders[j].IsBase &&
					(itfProviderHolders[i].Order > itfProviderHolders[j].Order ||
						itfProviderHolders[i].Provider.Name() > itfProviderHolders[j].Provider.Name())
		})
		for i, size := 0, len(itfProviderHolders); i < size-1; i++ {
			itfProviderHolders[i].InjectedIface, itfProviderHolders[i+1].InjectedIface =
				itfProviderHolders[i+1].InjectedIface, itfProviderHolders[i].InjectedIface
		}
		result.Bindings = append(result.Bindings, itfProviderHolders...)
	}

	return
}
