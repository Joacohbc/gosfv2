import PropTypes from 'prop-types';
import "./Button.css";

const Button = (props) => {
    return <button className={ props.className + " button"} onClick={props.onClick}>{props.text}</button>;
}

Button.propTypes = {
    className: PropTypes.string,
    onClick: PropTypes.func,
    text: PropTypes.any,
}

export default Button;