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
	Bindings  []*ProviderBinding
}

func (r *ParsedWireResult) HasProvider() bool {
	return len(r.Providers) > 0
}

type ProviderBinding struct {
	OriginIface   types.Type
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

	providerIfaceToProviders := map[types.Type][]types.Object{}
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
			providerIfaceToProviders[iface] = append(providerIfaceToProviders[iface], function)
		}

		firstParam := f.pkgParser.FirstParam(function)
		if firstParam == nil {
			continue
		}
	}

	for providerIface, providers := range providerIfaceToProviders {
		providerBindings := make([]*ProviderBinding, 0, len(providers))
		for _, provider := range providers {
			providerObject := f.pkgParser.FirstResult(provider)
			params := f.pkgParser.Params(provider)
			needInjectParam := false
			for _, param := range params {
				if !f.pkgParser.AssignableTo(param.Type(), providerIface) {
					continue
				}
				needInjectParam = true
				wireMeta := f.metaParser.ObjectMetaGroup(provider, MetaWireProvider)[0]
				providerHolder := &ProviderBinding{
					OriginIface:   providerIface,
					Provider:      provider,
					ProviderType:  providerObject.Type(),
					InjectedIface: param.Type(),
					Order:         Order(wireMeta),
				}
				providerBindings = append(providerBindings, providerHolder)
			}
			if !needInjectParam {
				providerBinding := &ProviderBinding{
					OriginIface:   providerIface,
					Provider:      provider,
					ProviderType:  providerObject.Type(),
					InjectedIface: providerIface,
					IsBase:        true,
				}
				providerBindings = append(providerBindings, providerBinding)
			}
		}

		f.sortProviderBindings(providerBindings)

		for i, size := 0, len(providerBindings); i < size-1; i++ {
			providerBindings[i].InjectedIface, providerBindings[i+1].InjectedIface =
				providerBindings[i+1].InjectedIface, providerBindings[i].InjectedIface
		}
		result.Bindings = append(result.Bindings, providerBindings...)
	}

	return
}

func (f *functions) sortProviderBindings(itfProviderBindings []*ProviderBinding) {
	sort.Slice(itfProviderBindings, func(i, j int) bool {
		return itfProviderBindings[i].IsBase ||
			!itfProviderBindings[j].IsBase &&
				(itfProviderBindings[i].Order > itfProviderBindings[j].Order ||
					itfProviderBindings[i].Provider.Name() > itfProviderBindings[j].Provider.Name())
	})
}
