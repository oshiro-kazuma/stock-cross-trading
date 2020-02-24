package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sclevine/agouti"
)

const (
	loginPage      = "https://site1.sbisec.co.jp/ETGate/?_ControlID=WPLETlgR001Control"
	salePage       = `https://site1.sbisec.co.jp/ETGate/?_PageID=WPLETstT001Mord10&_DataStoreID=DSWPLETstT001Control&_ControlID=WPLETstT001Control&_ActionID=order&trade_kbn=3&stock_sec_code=`
	buyPage        = `https://site1.sbisec.co.jp/ETGate/?_ControlID=WPLETstT002Control&_DataStoreID=DSWPLETstT002Control&stock_sec_code=`
	checkOrderPage = `https://site1.sbisec.co.jp/ETGate/?_ControlID=WPLETstT013Control&_PageID=DefaultPID&_DataStoreID=DSWPLETstT013Control&_SeqNo=1582514728899_default_task_609_WPLETstT014Cplc10_list&getFlg=on&_ActionID=DefaultAID`
)

var (
	loginID         = "fixme" // ログインID
	loginPass       = "fixme" // ログインパスワード
	transactionPass = "fixme" // 取引パスワード
	orders          = []Order{
		{
			Code:     "2670", // 注文を入れたい銘柄コード
			Quantity: "100",  // 注文を入れる株数
		},
	}
	dryRun = false
)

var (
	driver = agouti.ChromeDriver()
	page   *agouti.Page
)

func init() {
	driver.Debug = true

	if err := driver.Start(); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start driver"))
	}

	_page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		log.Fatal(err)
	}

	page = _page
}

type Order struct {
	Code     string
	Quantity string
}

func main() {
	// 終了時にchromeを起動したままにする
	defer driver.Stop()

	// ログイン
	err := login()
	if err != nil {
		log.Fatal(err)
	}

	// 信用売り注文
	for _, order := range orders {
		err = sale(order)
		if err != nil {
			fileName := fmt.Sprintf("/tmp/stock_order_sale_failed_%d_%s_%s", time.Now().Unix(), order.Code, order.Quantity)
			page.Screenshot(fileName)

			fmt.Println("order sale failed. screenshot: " + fileName)
		} else {
			fileName := fmt.Sprintf("/tmp/stock_order_sale_success_%d_%s_%s", time.Now().Unix(), order.Code, order.Quantity)
			page.Screenshot(fileName)

			fmt.Println("order sale success. screenshot: " + fileName)
		}
	}

	// 現物買い注文
	for _, order := range orders {
		err = buy(order)
		if err != nil {
			fileName := fmt.Sprintf("/tmp/stock_order_buy_falied_%d_%s_%s_failed", time.Now().Unix(), order.Code, order.Quantity)
			page.Screenshot(fileName)

			fmt.Println("order buy failed. screenshot: " + fileName)
		} else {
			fileName := fmt.Sprintf("/tmp/stock_order_buy_success_%d_%s_%s", time.Now().Unix(), order.Code, order.Quantity)
			page.Screenshot(fileName)

			fmt.Println("order buy success. screenshot: " + fileName)
		}
	}

	// 注文照会ページに遷移
	page.Navigate(checkOrderPage)

	time.Sleep(10 * time.Minute)
}

func login() error {
	// navigate login page
	err := page.Navigate(loginPage)
	if err != nil {
		return err
	}

	// send login account
	err = page.FindByXPath(`//input[@name='user_id']`).SendKeys(loginID)
	if err != nil {
		return err
	}
	err = page.FindByXPath(`//input[@name='user_password']`).SendKeys(loginPass)
	if err != nil {
		return err
	}
	err = page.FindByXPath(`//input[@name='logon']`).Click()
	if err != nil {
		return err
	}

	return nil
}

func buy(order Order) error {
	// 注文ページに遷移
	err := page.Navigate(buyPage + order.Code)
	if err != nil {
		return err
	}
	// SOR指定のチェックを外す
	err = checkOff(page.FindByXPath(`//*[@id="sor_check"]`))
	if err != nil {
		return err
	}
	// 株数
	err = page.FindByXPath(`//*[@id="input_quantity"]`).SendKeys(order.Quantity)
	if err != nil {
		return err
	}
	// 成行き
	err = page.FindByXPath(`//*[@id="gsn1"]/label`).Click()
	if err != nil {
		return err
	}
	// 取引パスワード
	err = page.FindByXPath(`//*[@id="pwd3"]`).SendKeys(transactionPass)
	if err != nil {
		return err
	}
	// 注文確認画面を省略
	err = checkOn(page.FindByXPath(`//*[@id="shouryaku"]`))
	if err != nil {
		return err
	}
	// 注文実行
	if !dryRun {
		err = page.FindByXPath(`//*[@id="botton2"]/img`).Click()
		if err != nil {
			return err
		}
	} else {
		time.Sleep(3 * time.Second)
	}

	return nil
}

// 信用売りを入れる
func sale(order Order) error {
	// 注文ページに遷移
	err := page.Navigate(salePage + order.Code)
	if err != nil {
		return err
	}
	// SOR指定のチェックを外す
	err = checkOff(page.FindByXPath(`//*[@id="sor_check"]`))
	if err != nil {
		return err
	}
	// 株数
	err = page.FindByXPath(`//*[@id="input_quantity"]`).SendKeys(order.Quantity)
	if err != nil {
		return err
	}
	// 成行き
	err = page.FindByXPath(`//*[@id="gsn1"]/label`).Click()
	if err != nil {
		return err
	}
	// 一般信用(15日)
	err = page.FindByXPath(`//*[@id="payment_limit_labelA"]`).Click()
	if err != nil {
		return err
	}
	// 取引パスワード
	err = page.FindByXPath(`//*[@id="pwd3"]`).SendKeys(transactionPass)
	if err != nil {
		return err
	}
	// 注文確認画面を省略
	err = checkOn(page.FindByXPath(`//*[@id="shouryaku"]`))
	if err != nil {
		return err
	}
	// 注文実行
	if !dryRun {
		err = page.FindByXPath(`//*[@id="botton2"]/img`).Click()
		if err != nil {
			return err
		}
	} else {
		time.Sleep(3 * time.Second)
	}

	return nil
}

// チェックボックスをOFF状態にする
func checkOff(selection *agouti.Selection) error {
	if selection == nil {
		return fmt.Errorf("shouryaku checkbox is not found")
	}
	selected, err := selection.Selected()
	if err != nil {
		return err
	}
	if selected {
		if err := selection.Click(); err != nil {
			return err
		}
	}
	return nil
}

// チェックボックスをON状態にする
func checkOn(selection *agouti.Selection) error {
	if selection == nil {
		return fmt.Errorf("shouryaku checkbox is not found")
	}
	selected, err := selection.Selected()
	if err != nil {
		return err
	}
	if !selected {
		if err := selection.Click(); err != nil {
			return err
		}
	}
	return nil
}
