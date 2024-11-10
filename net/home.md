# home network layout

## current (Nov 2024)

```mermaid
architecture-beta
    group house(house)[House]
    group pb(house)[Pole Barn]

    service lr(wifi)[living room AP] in house
    service base(router)[Router] in house
    service apt(wifi)[Apartment] in house
    service ape(wifi)[Airport Extreme] in pb
    service plu(plug)[Powerline (utility room)] in house
    service pla(plug)[Powerline (attic)]
    service wt(iot)[Wireless Tags Base Station]
    service ur(router)[Router (landing)]
    service matt(computer)[Matt's Desk]
    service planes(iot)[ADSB receiver]
    service oliver(computer)[Oliver's desk]
    service tiny(server)[Curly-goggles]

    lr:R -- L:base
    base:R -- L:apt
    base:B -- T:ape
    base:T -- B:plu
    plu:R -- L:pla
    pla:R -- L:wt
    base:T -- R:ur
    ur:T -- B:matt
```

## proposed