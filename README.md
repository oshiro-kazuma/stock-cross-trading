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

約定の指定なしで注文を入れているので、市場が開いていないときに実行すること。

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

## appendix

### クロス取引可能な銘柄の探し方

下記のページから検索可能

https://trading0.sbisec.co.jp/cgss/yutai/yutai_search.do#sort=PS_CODE-ASC&menuChange=true

### 空売り在庫の調べ方

下記ページから検索可能。「優待権利確定月」を指定して、「優待権利確定月」でソートすると探しやすい。

https://site0.sbisec.co.jp/marble/domestic/top/cbsProductList.do

### 解説動画

https://www.youtube.com/watch?v=2aT_P2Pea3A ここ見てやりたくなってしまった。

### FAQ
 - 逆日歩は大丈夫なの？
   - 一般信用取引を使えば問題ない
 - 片方だけ約定したりしないの？(売りだけ入るとか買いだけ入っちゃうとか)
   - ポジションを取るときは、相場が開く前に買いと売りで両方の成行注文を入れているため大丈夫っぽい（おそらく必ず寄り付く？）
   - ポジション解決は「現渡」でやるので大丈夫
 - 儲かるの？
   - ゆーてタネ銭必要だし、そんな儲からないので自動化して回したい。
   - オペミスしなければノーリスクで優待取得が可能
 - いつポジションを取るのがいいの？
   - 権利確定日の9時までに注文を入れておけばいいが、人気な銘柄は一般信用の在庫がなくなるので「月末の権利付最終売買日から14営業日前」にやるのが良さそう。
   - https://faq.sbisec.co.jp/faq_detail.html?id=46510&category=67&page=1 ここ見たほうがいい。
