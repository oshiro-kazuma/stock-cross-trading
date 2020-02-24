# stock-cross-trading

## これはなに

SeleniumでSBI証券にログインして、優待のタダ取り注文であるクロス取引の注文を入れるボットです。

人気銘柄は空売り在庫がすぐなくなるらしいので、ボットで購入できるように。

スピード注文を意識はしていないので巷のボットには負けるかも。

次の3月権利日15日前に試してみる。

## 注文仕様

在庫を先に取りたいのので空売りの注文を先に入れ、次に現物で買い注文を入れています。

1. 信用新規売 （一般信用/15日）の成行を入れる。
2. 現物買の成行きを入れる。市場は東証でSOR指定のチェックを外す。
3. 権利落ち日まで待つ。
4. 現渡での手じまいは未実装のため手動で行う。

## setup for mac

selenium関連のインストール

```sh
brew install selenium-server-standalone
brew cask install chromedriver
```

## install

```sh
git clone git@github.com:oshiro-kazuma/stock-cross-trading.git
cd stock-cross-trading
make install
xtrade
```

## 動かし方

パスワード等を引数化するのが面倒だったのでハードコーディング中。

main.goの下記を修正する

 - loginID
 - loginPass
 - transactionPass
 - orders
 
 その後、install手順を実行して出来上がった実行バイナリである `xtrade` を実行する。
