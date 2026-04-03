# Architecture Decisions

## Mode Split

Aeroform uses one CLI with two operating modes instead of separate tools.

## Constrained Generation

The AI layer selects from validated templates rather than inventing arbitrary infrastructure.

## Provider Abstraction

Cloud-specific behavior is isolated behind provider interfaces so the engine can stay shared.

## Security Model

Security checks are mode-aware: beginner-friendly in Simple Mode, strict and technical in Pro Mode.

