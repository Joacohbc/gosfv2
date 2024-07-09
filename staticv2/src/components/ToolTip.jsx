import 'bootstrap/dist/css/bootstrap.min.css';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import BoostrapTooltip from 'react-bootstrap/Tooltip';
import PropTypes from 'prop-types';

const ToolTip = ({ children, toolTipMessage, placement }) => (
    <OverlayTrigger 
        placement={placement}
        overlay={<BoostrapTooltip>{toolTipMessage}</BoostrapTooltip>}
        delay={100}>
        {children}
    </OverlayTrigger>
);

ToolTip.propTypes = {
    children: PropTypes.element.isRequired,
    toolTipMessage: PropTypes.string.isRequired,
    placement: PropTypes.oneOf(['top', 'bottom', 'left', 'right'])
};

export default ToolTip;