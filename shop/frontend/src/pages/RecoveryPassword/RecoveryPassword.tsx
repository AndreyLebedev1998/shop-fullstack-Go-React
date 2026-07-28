import { useState } from "react";
import {
  Alert,
  Button,
  Form,
  InputGroup,
  Modal,
  OverlayTrigger,
  Tooltip
} from "react-bootstrap";
import { useTranslation } from "react-i18next";
import "./recovery-password.css";
import { codeMatching, sendCodeEmail, sendCodeTg } from "../../api/auth-api/auth-api";
import { useDispatch } from "react-redux";
import { setEmailRecoveryPassword, setNewToken } from "../../store/authSlice";
import { useNavigate } from "react-router-dom";

const RecoveryPassword = () => {
  const { t } = useTranslation();
  const [email, setEmail] = useState<string>("");
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [successMessage, setSuccessMeaage] = useState(false);
  const [showSuccessMessage, setShowSuccessMessage] = useState(false);
  const [errorMessage, setErrorMessage] = useState(false);
  const [showError, setShowError] = useState(false);
  const [errorMessageChatId, setErrorMessageChatId] = useState(false);
  const [showErrorMessageChatId, setShowErrorMessageChatId] = useState(false);
  const [showModalRecoveryPassword, setShowRecoveryPassword] = useState(false);
  const [errorCodeMathing, setErrorCodeMathing] = useState(false)
  const [showErrorCodeMatching, setShowErrorCodeMatching] = useState(false);
  const [errorEmailIsNoDefined, setErrorEmailIsNoDefined] = useState(false);
  const [errorCodeMathingEmail, setErrorCodeMathingEmail] = useState(false);
  const [showSuccessMessageEmail, setShowSuccessMessageEmail] = useState(false); 
  const [disabled, setDisabled] = useState(false)
  const [code, setCode] = useState("");

  const handleSendTg = () => {
    if (!email && !email.includes("@")) return;
    setDisabled(true)
    sendCodeTg(
      email,
      setErrorMessage,
      setShowError,
      setShowErrorMessageChatId,
      setErrorMessageChatId,
      setSuccessMeaage
    ).then((data) => {
      if (data) {
        setSuccessMeaage(true);
        setShowSuccessMessage(true);
        setShowRecoveryPassword(true);
        setErrorMessage(false)
      }
    }).finally(() => {
      setDisabled(false)
    })
  };

  const handleSendEmail = () => {
    if (!email && !email.includes("@")) return;
    setDisabled(true)
    sendCodeEmail(
      email,
      setErrorEmailIsNoDefined,
      setErrorCodeMathingEmail,
    ).then((data) => {
      if (data) {
        setSuccessMeaage(true);
        setErrorEmailIsNoDefined(false);
        setShowSuccessMessageEmail(true);
        setShowRecoveryPassword(true);

      }
    }).finally(() => {
      setDisabled(false)
    })
  };

  const handleCodeMatching = () => {
    if (code.length != 6) return;
    codeMatching({email, phone: "", code,}, setErrorCodeMathing, setShowErrorCodeMatching).then((data) => {
      if (data) {
        dispatch(setNewToken({ token: data.token }))
        setTimeout(() => {
          dispatch(setEmailRecoveryPassword(email))
          navigate("/new-password")
        }, 100)
      }
    })
  }

  function submit() {}
  return (
    <>
      {showModalRecoveryPassword && (
        <Modal
          size="lg"
          aria-labelledby="contained-modal-title-vcenter"
          centered
          show={showModalRecoveryPassword}
          onHide={() => {
            setShowRecoveryPassword(false)
            setErrorCodeMathing(false)
          }}
        >
          <Modal.Header closeButton>
            <Modal.Title id="contained-modal-title-vcenter">
              {t("auth.enter_the_code")}
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            {errorCodeMathing && showErrorCodeMatching && (
            <Alert
              className="code-matching-alert"
              onClose={() => setShowErrorCodeMatching(false)}
              dismissible
              variant="danger"
            >
              {t("auth.error")}
            </Alert>
          )}
            <Form.Control type="number" onChange={(e) => setCode(e.target.value)} aria-describedby="passwordHelpBlock"/>
              <Form.Text id="passwordHelpBlock" muted>{t("auth.enter_the_code_six_symbols")}</Form.Text>
          </Modal.Body>
          <Modal.Footer>
            <Button onClick={() => {
              handleCodeMatching()
            }}>
              {t("auth.send")}
            </Button>
          </Modal.Footer>
        </Modal>
      )}
      <div className="recovery-password-page">
        <div className="recovery-password-wrap">
          {errorMessage && showError && (
            <Alert
              className="send-code-alert"
              onClose={() => setShowError(false)}
              dismissible
              variant="danger"
            >
              {t("auth.error")}
            </Alert>
          )}
          {errorEmailIsNoDefined && (
            <Alert
              className="error-alert-email"
              onClose={() => setErrorEmailIsNoDefined(false)}
              dismissible
              variant="danger"
            >
              {t("auth.invalid_email")}
            </Alert>
          )}
          {errorCodeMathingEmail && (
            <Alert
              className="error-alert-email"
              onClose={() => setErrorCodeMathingEmail(false)}
              dismissible
              variant="danger"
            >
              {t("auth.error_send_code_email")}
            </Alert>
          )}
          {successMessage && showSuccessMessage && (
            <>
              <Alert
                className="message-success-send-code"
                onClose={() => setShowSuccessMessage(false)}
                dismissible
                variant="success"
              >
                {t("auth.code_sent_successfully")}
              </Alert>
              <span>{t("auth.expectation_one_minute")}</span>
            </>
          )}
          {showSuccessMessageEmail && (
            <>
              <Alert
                className="message-success-send-code"
                onClose={() => setShowSuccessMessage(false)}
                dismissible
                variant="success"
              >
                {t("auth.code_sent_successfully")}
              </Alert>
              <span>{t("auth.expectation_one_minute")}</span>
            </>
          )}
          {errorMessageChatId && showErrorMessageChatId && (
            <>
              <Alert
                className="message-not-chat-id"
                onClose={() => setShowErrorMessageChatId(true)}
                dismissible
                variant="danger"
              >
                {t("auth.not_linked_telegram")}
              </Alert>
            </>
          )}
          <InputGroup className="mb-3">
            <InputGroup.Text id="basic-addon1">@</InputGroup.Text>
            <Form.Control
              name="email"
              autoComplete="email"
              placeholder={t("auth.email")}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  submit();
                }
              }}
            />
          </InputGroup>
          <div className="recovery-password-actions">
            <OverlayTrigger
              placement="top"
              overlay={
                <Tooltip>{t("auth.only_if_telegram_is_linked")}</Tooltip>
              }
            >
              <Button
                className="send-code-tg-button"
                type="button"
                variant="info"
                onClick={handleSendTg}
                disabled={disabled}
              >
                {t("auth.send_code_telegram")}
                <i className="bi bi-telegram send-code-tg-icon"></i>
              </Button>
            </OverlayTrigger>
            <Button
              className="send-code-tg-email"
              type="button"
              variant="success"
              onClick={handleSendEmail}
              disabled={disabled}
            >
              {t("auth.send_code_email")}{" "}
              <i className="bi bi-envelope send-code-email-icon"></i>
            </Button>
          </div>
        </div>
      </div>
    </>
  );
};

export default RecoveryPassword;
