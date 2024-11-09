# home network layout

## current (Nov 2024)

```mermaid
architecture-beta
    group house(house)[House]

    service lr(wifi)[living room AP] in house
    service base(wifi)[Router] in house
    service apt(wifi)[Apartment] in house

    lr:R -- L:base
    base:R -- L:apt
```

## proposed