(let arr (mk-array))
(append arr 3 4 5)
(let i 1)
(nth i (append arr 7))
(nth 2 (append (mk-array) 3 (+ 2 3) 7))
(fn car (arr) (nth 0 arr))
(car arr)
