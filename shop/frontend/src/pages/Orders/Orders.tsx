import { useEffect, useState, type FC } from "react";
import type { RootState } from "../../store/store";
import { useSelector } from "react-redux";
import { useTranslation } from "react-i18next";
import "./orders.css";
import { getOrdersAllTime } from "../../api/orders-api/orders-api";
import { Accordion, Alert } from "react-bootstrap";
import type { OrderType } from "../../../types/types";
import { config } from "../../config";

const Orders: FC = () => {
  const token = useSelector((state: RootState) => state.auth.token);
  const [orders, setOrders] = useState<OrderType[]>([]);
  const [error, setError] = useState(false);
  const [show, setShow] = useState(false);
  const { t } = useTranslation();

  useEffect(() => {
    if (token) {
      getOrdersAllTime(token, setError, setShow).then((data) => {
        if (data) {
          setOrders(data);
        }
      });
    }
  }, [token]);
  return (
    <div className="orders-page">
      {error && show && (
        <Alert onClose={() => setShow(false)} dismissible variant="danger">
          {t("orders.error_receiving_orders")}
        </Alert>
      )}
      {token ? (
        <div className="orders-page__orders-wrapper">
          <Accordion alwaysOpen>
            {orders.map((order) => {
              return (
                <Accordion.Item
                  key={order.order_id}
                  eventKey={String(order.order_id)}
                >
                  <Accordion.Header>
                    <div className="order-header">
                      <span>{t("orders.order")} №{order.order_id}</span>
                      <span>
                        {new Date(order.created_at).toLocaleString("ru-RU", {
                          day: "2-digit",
                          month: "2-digit",
                          year: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </span>
                    </div>
                  </Accordion.Header>

                  <Accordion.Body>
                    <div className="order-info">
                      <div>
                        <strong>{t("orders.email")}:</strong> {order.email}
                      </div>

                      <div>
                        <strong>{t("orders.phone")}:</strong> {order.phone}
                      </div>

                      <div>
                        <strong>{t("orders.status")}:</strong>{" "}
                        {t(`orders.status_type.${order.status}`)}
                      </div>

                      <div>
                        <strong>{t("orders.total_amount")}:</strong>{" "}
                        {order.total_price} ₽
                      </div>
                    </div>

                    <div className="order-products">
                      {order.products?.map((product) => (
                        <div key={product.product_id} className="order-product">
                          <img
                            className="image-order-product"
                            src={`${config.PRODUCTS_IMAGES_BASE_URL}/${product.image_url}`}
                            alt={product.product_name}
                          />

                          <div className="order-product-info">
                            <div className="order-product-name">
                              {product.product_name}
                            </div>

                            <div>
                              {t("orders.quantity")}: {product.quantity}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </Accordion.Body>
                </Accordion.Item>
              );
            })}
          </Accordion>
        </div>
      ) : (
        <div className="orders-warn">{t("orders.warn_orders")}</div>
      )}
    </div>
  );
};

export default Orders;
