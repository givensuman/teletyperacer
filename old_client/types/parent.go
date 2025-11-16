package types

import "github.com/charmbracelet/bubbletea"

// SingleParent defines a tea.Model with only one child.
type SingleParent[T tea.Model] struct {
	Child T
}

// Replace replaces a single parent's child.
func (sp *SingleParent[T]) Replace(child T) {
	sp.Child = child
}

type ComparableModel interface {
	comparable
	tea.Model
}

// Parent defines a tea.Model with multiple children.
type Parent[T ComparableModel] struct {
	Children []T
}

// Add adds a child to the parent's state.
func (p *Parent[T]) Add(child T) *Parent[T] {
	p.Children = append(p.Children, child)
	return p
}

// Remove removes a child from the parent's state.
func (p *Parent[T]) Remove(child T) *Parent[T] {
	for i := range p.Children {
		if child == p.Children[i] {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
		}
	}

	return p
}

// UpdateChildren calls each child's Update function.
func (p Parent[T]) UpdateChildren(msg tea.Msg) {
	for _, child := range p.Children {
		child.Update(msg)
	}
}
