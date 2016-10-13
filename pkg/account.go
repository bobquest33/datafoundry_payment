package pkg

import (
	"net/http"

	"github.com/zonesan/clog"
)

type AccountAgent service

type Account struct {
	Purchased bool    `json:"purchased"`
	Notify    bool    `json:"notification"`
	Plans     []Plan  `json:"subscriptions,omitempty"`
	Status    string  `json:"status"`
	Balance   Balance `json:"balance"`
}

func (agent *AccountAgent) Get(r *http.Request) (*Account, error) {
	r.ParseForm()

	project := r.FormValue("project")

	clog.Debug(project)

	_, err := getToken(r)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	account := new(Account)

	if orders, err := agent.Checkout.ListOrders(r); err != nil {
		clog.Error(err)
		return nil, err
	} else {
		//clog.Debugf("%#v", orders)

		if len(*orders) > 0 {
			account.Purchased = true

			if plans, err := agent.Market.List(r); err != nil {
				clog.Error(err)
			} else {
				func() {
					for _, order := range *orders {
						found := false
						for _, plan := range plans.Plans {
							if order.Plan_id == plan.PlanId {
								account.Plans = append(account.Plans, plan)
								clog.Debug(order.Plan_id, "found in plan list.")
								found = true
								break
							}
						}
						if !found {
							clog.Warnf("order with plan id '%v' not found in market.", order.Plan_id)
						}
					}

				}()
			}

			// for _, order := range *orders {
			// 	if plan, err := agent.Market.Get(r, order.Plan_id); err != nil {
			// 		clog.Error(err)
			// 	} else {
			// 		if plan.PlanId != "" {
			// 			account.Plans = append(account.Plans, *plan)
			// 		} else {
			// 			clog.Warn("empty plan id ...", plan)
			// 		}
			// 	}
			// }
		}

	}

	//account := fakeAccount(r)
	return account, nil
}
