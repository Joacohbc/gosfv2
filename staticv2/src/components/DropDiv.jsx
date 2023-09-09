import { useEffect, useRef } from "react";

const DropDiv = ({ onDrop, children }) => {
    const drop = useRef(null);

    useEffect(() => {
        const div = drop.current;

        div.addEventListener("dragover", (e) => {
            e.preventDefault();
            e.stopPropagation();
            div.classList.add("dragover");
        });

        div.addEventListener("dragleave", (e) => {
            e.preventDefault();
            e.stopPropagation();
            div.classList.remove("dragover");
        });

        div.addEventListener("drop", (e) => {
            e.preventDefault();
            e.stopPropagation();
            div.classList.remove("dragover");

            const filesDropped = [];
            if (e.dataTransfer.items) {
                for (let i = 0; i < e.dataTransfer.items.length; i++) {
                    if (e.dataTransfer.items[i].kind === "file") {
                        filesDropped.push(e.dataTransfer.items[i].getAsFile());
                    }
                }
                e.dataTransfer.items.clear();
            } else {
                for (let i = 0; i < e.dataTransfer.files.length; i++) {
                    filesDropped.push(e.dataTransfer.files[i]);
                }
                e.dataTransfer.clearData();
            }

            onDrop(filesDropped);
        });
    }, [drop, onDrop]);

    return <div ref={drop}> {children} </div>;
};

export default DropDiv;
