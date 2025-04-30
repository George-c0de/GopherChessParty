import styled from 'styled-components';

export const Container = styled.div`
    max-width: 400px;
    margin: 50px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
`;

export const FormGroup = styled.div`
    margin-bottom: 15px;
`;

export const Label = styled.label`
    display: block;
    margin-bottom: 5px;
`;

export const Input = styled.input`
    width: 100%;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
`;

export const Button = styled.button`
    background-color: #4CAF50;
    color: white;
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    width: 100%;

    &:hover {
        background-color: #45a049;
    }
`;

export const ErrorMessage = styled.div`
    color: red;
    margin-top: 10px;
    display: ${props => props.hidden ? 'none' : 'block'};
`;

export const Link = styled.a`
    color: #4CAF50;
    text-decoration: none;

    &:hover {
        text-decoration: underline;
    }
`;

export const LinkContainer = styled.div`
    text-align: center;
    margin-top: 15px;
`; 