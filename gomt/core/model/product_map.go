package model

import (
	"fmt"
	"strings"
	"sync"
)

type ProductMapNode struct {
	Level    string
	Code     string
	Name     string
	Children []*ProductMapNode
}

func (n ProductMapNode) Dump(prefix string) string {
	tmp := fmt.Sprintf("%v%v: %v, %v\n", prefix, n.Level, n.Code, n.Name)
	for _, v := range n.Children {
		tmp += v.Dump(prefix + "-")
	}
	return tmp
}

var defaultProductMapRoot *ProductMapNode
var defaultProductMapRootOnce sync.Once

func GetDefaultProductMapRoot() *ProductMapNode {
	defaultProductMapRootOnce.Do(func() {
		root := createDefaultProductMapRoot()
		defaultProductMapRoot = root
	})
	return defaultProductMapRoot
}

func createDefaultProductMapRoot() *ProductMapNode {
	root := ProductMapNode{"root", "root", "root", []*ProductMapNode{}}
	maps := GetDefaultCodingDefineMap()

	for familyId, familyName := range maps.FamilyMap {
		familyNode := ProductMapNode{"ProductFamily", familyId, familyName, []*ProductMapNode{}}
		for productTypeId, productTypeName := range maps.ProductTypeMap {
			if strings.HasPrefix(productTypeId, familyId) {
				productNode := ProductMapNode{"ProductType", productTypeId, productTypeName, []*ProductMapNode{}}
				for platformId, platformName := range maps.PlatformMap {
					if strings.HasPrefix(platformId, productTypeId) {
						platformNode := ProductMapNode{"Platform", platformId, platformName, []*ProductMapNode{}}
						for revisionId, revisionName := range maps.RevisionMap {
							if strings.HasPrefix(revisionId, platformId) {
								revisionNode := ProductMapNode{"Revision", revisionId, revisionName, []*ProductMapNode{}}
								platformNode.Children = append(platformNode.Children, &revisionNode)
							}
						}
						productNode.Children = append(productNode.Children, &platformNode)
					}
				}
				familyNode.Children = append(familyNode.Children, &productNode)
			}
		}
		root.Children = append(root.Children, &familyNode)
	}
	return &root
}
