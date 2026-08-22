package tools

import (
	"context"

	"github.com/jamesawo/mdev/internal/infrastructure/environment"
)

// ContextTool is an optional extension implemented by tools whose active
// lifecycle operations can be cancelled.
type ContextTool interface {
	InstallContext(context.Context, *environment.Environment) error
	ConfigureContext(context.Context, *environment.Environment) error
	VerifyContext(context.Context, *environment.Environment) error
}

// InstallContext invokes the context-aware install operation when available.
func InstallContext(ctx context.Context, tool Tool, env *environment.Environment) error {
	if contextual, ok := tool.(ContextTool); ok {
		return contextual.InstallContext(ctx, env)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return tool.Install(env)
}

// ConfigureContext invokes the context-aware configure operation when available.
func ConfigureContext(ctx context.Context, tool Tool, env *environment.Environment) error {
	if contextual, ok := tool.(ContextTool); ok {
		return contextual.ConfigureContext(ctx, env)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return tool.Configure(env)
}

// VerifyContext invokes the context-aware verification operation when available.
func VerifyContext(ctx context.Context, tool Tool, env *environment.Environment) error {
	if contextual, ok := tool.(ContextTool); ok {
		return contextual.VerifyContext(ctx, env)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return tool.Verify(env)
}
