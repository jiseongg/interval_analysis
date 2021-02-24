define i32 @f(i32 %a) {
entry:
  %b = add i32 0, 1
  %c = add i32 0, 2
  %d = add i32 %c, 3
  %e = add i32 %d, %c
  ret i32 %e
}

