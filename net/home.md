# home network layout

## current (Nov 2024)

```mermaid
architecture-beta
    group house(house)[House]
    group pb(house)[Pole Barn]

    service lr(wifi)[living room AP] in house
    service base(wifi)[Router] in house
    service apt(wifi)[Apartment] in house
    service ape(wifi)[Airport Extreme] in pb
    service plu(wifi)[Powerline (utility room)] in house
    service pla(wifi)[Powerline (attic)]
    service wt(iot)[

    lr:R -- L:base
    base:R -- L:apt
    base:B -- T:ape
```

## proposed