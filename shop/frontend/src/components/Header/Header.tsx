import { useEffect, useState, type FC } from "react";
import {
  Container,
  Nav,
  Navbar,
  NavDropdown,
  Offcanvas,
} from "react-bootstrap";
import RuFlag from "country-flag-icons/react/3x2/RU";
import GbFlag from "country-flag-icons/react/3x2/GB";
import Logo from "../../images/ready_fresh_mart_logo.png";
import type { LinksType } from "../../../types/types.ts";
import { getCategoriesWithSubcategories } from "../../api/products-api/products-api.ts";
import { useTranslation } from "react-i18next";
import "./header.css";
import { Link } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../../store/store.ts";
import { logout } from "../../store/authSlice.ts";
import { setSubcategories } from "../../store/categoriesSlice.ts";

interface HeaderType {
  setChoceSubcategory: (num: number) => void
}

const Header: FC<HeaderType> = ({setChoceSubcategory}) => {
  const { t, i18n } = useTranslation();
  const dispatch = useDispatch();
  const token = useSelector((state: RootState) => state.auth.token);
  const categories = useSelector((state: RootState) => state.categories.categories)
  const [show, setShow] = useState(false);
  const links: LinksType[] = [
    { link: "/", name: t("header.links.home"), active: true },
    { link: "/orders", name: t("header.links.orders"), active: true },
    { link: "/about-us", name: t("header.links.about_us"), active: true },
    { link: "/profile", name: t("header.links.profile"), active: !!token },
  ];

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);

  useEffect(() => {
    getCategoriesWithSubcategories().then((data) =>
      dispatch(setSubcategories(data))
    );
  }, [dispatch]);

  const exit = () => {
    dispatch(logout());
  };

  return (
    <>
      <Offcanvas show={show} onHide={handleClose}>
        <Offcanvas.Header closeButton>
          <Offcanvas.Title>Меню</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
            {categories != null && categories.map((category) => (
                <div key={category.id} className="d-flex align-items-center">
                  <Link
                    to={`/products-categories/${category.id}`}
                    className="flex-grow-1 py-2 px-3 text-decoration-none text-dark"
                    onClick={() => {
                      handleClose()
                      setChoceSubcategory(0)
                    }}
                  >
                    {t(`header.offcanvas.categories.${category.category_name}`)}
                  </Link>
                </div>
            ))}
        </Offcanvas.Body>
      </Offcanvas>
      <Navbar expand="lg" className="bg-body-tertiary header__navbar">
        <Container className="header__container">
          <Navbar.Toggle
            aria-controls="basic-navbar-nav"
            className="header__button-burger"
            onClick={handleShow}
          />
          <Navbar.Brand className="header__brand">
            <Link to={"/"}>
            <img className="header__img-logo" src={Logo} />
            </Link>
          </Navbar.Brand>
          <Navbar.Collapse
            className="header__link-wrapper"
            id="basic-navbar-nav"
          >
            {links.map(
              (link) =>
                link.active && (
                  <Link key={link.link} to={link.link}>
                    {link.name}
                  </Link>
                ),
            )}
          </Navbar.Collapse>
          <Nav className="me-auto">
            <NavDropdown
              drop="start"
              className="header__dropdown-change-language"
              title={<i className="header_bi-globe-icon bi bi-globe"></i>}
            >
              <NavDropdown.Item
                onClick={() => {
                  i18n.changeLanguage("ru");
                  localStorage.setItem("language", "ru");
                }}
                active={i18n.language === "ru"}
              >
                <RuFlag style={{ width: "20px", marginRight: "8px" }} /> RU
              </NavDropdown.Item>
              <NavDropdown.Item
                onClick={() => {
                  i18n.changeLanguage("en");
                  localStorage.setItem("language", "en");
                }}
                active={i18n.language === "en"}
              >
                <GbFlag style={{ width: "20px", marginRight: "8px" }} /> EN
              </NavDropdown.Item>
            </NavDropdown>
          </Nav>
          {token ? (
            <i onClick={exit} className="bi bi-box-arrow-left"></i>
          ) : (
            <Link to={"/authorization"}>
              <i className="header_bi-person-icon bi bi-person"></i>
            </Link>
          )}
        </Container>
      </Navbar>
    </>
  );
};

export default Header;
