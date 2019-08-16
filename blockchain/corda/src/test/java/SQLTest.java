
import java.util.List;

import org.apache.commons.lang.StringUtils;
import org.junit.Test;

import net.sf.jsqlparser.JSQLParserException;
import net.sf.jsqlparser.expression.BinaryExpression;
import net.sf.jsqlparser.expression.Expression;
import net.sf.jsqlparser.expression.ExpressionVisitorAdapter;
import net.sf.jsqlparser.expression.LongValue;
import net.sf.jsqlparser.expression.Parenthesis;
import net.sf.jsqlparser.expression.operators.conditional.AndExpression;
import net.sf.jsqlparser.expression.operators.conditional.OrExpression;
import net.sf.jsqlparser.parser.CCJSqlParserUtil;
import net.sf.jsqlparser.statement.*;
import net.sf.jsqlparser.statement.select.*;
import net.sf.jsqlparser.statement.values.ValuesStatement;
import net.sf.jsqlparser.util.TablesNamesFinder;
import net.sf.jsqlparser.util.deparser.ExpressionDeParser;
import net.sf.jsqlparser.statement.values.*;
import net.sf.jsqlparser.expression.operators.relational.*;
import net.sf.jsqlparser.expression.*;

public class SQLTest {
	@Test
	public void testsql() throws JSQLParserException {
		Statement stmt = CCJSqlParserUtil.parse("select * from com.example.iou.IOU where linearId.externalId='abc' and ");
		Select select = (Select)stmt;
		
		TablesNamesFinder tablesNamesFinder = new TablesNamesFinder();
		List<String> tableList = tablesNamesFinder.getTableList(select);
		
		System.out.println(tableList.toString());
		
		select.getSelectBody().accept(new MySelectVisitor());
		
	}
	
	@Test
	public void testwhere() throws JSQLParserException {
		Expression expr = CCJSqlParserUtil.parseCondExpression("quantity between 2 and 100 ");
		System.out.println(expr.getClass().getName());
		System.out.println(((Between)expr).getBetweenExpressionStart());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity >= 100 ");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity <= 100 ");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity = 2019-08-15T13:30:00Z");
		System.out.println(expr.getClass().getName());
		expr = ((EqualsTo)expr).getRightExpression();
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity != 100 ");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity in ( 100, 200)");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity not in ( 100, 200)");
		System.out.println(expr.getClass().getName());
		
		System.out.println(((InExpression)expr).getLeftExpression().toString());
		System.out.println(((InExpression)expr).isNot());
		System.out.println(((InExpression)expr).getRightItemsList().toString());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity >= 100 or quantity <=10");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity >= 100 and quantity <=10");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("quantity >= 100 or (quantity between 0 and 100)");
		System.out.println(expr.getClass().getName());
		
		expr = CCJSqlParserUtil.parseCondExpression("sum(quantity)");
		System.out.println(expr.getClass().getName());
		
		ExpressionDeParser parser = new ExpressionDeParser();
		expr.accept(parser);
		
		
	}
	
	static class MySelectVisitor implements SelectVisitor {

		@Override
		public void visit(PlainSelect arg0) {
			System.out.println("where=" + arg0.getWhere().toString());
			System.out.println("alias=" + arg0.getFromItem().getAlias().getName());
			System.out.println("from=" + arg0.getFromItem().toString());
			arg0.getSelectItems().forEach(it -> System.out.println(it.toString()));
			//long offset = ((LongValue)arg0.getLimit().getOffset()).getValue();
			
			arg0.getWhere().accept(new FilterExpressionVisitorAdapter() );
		}

		@Override
		public void visit(SetOperationList arg0) {
			System.out.println("SetOperationList=" + arg0.toString());
		}

		@Override
		public void visit(WithItem arg0) {
			System.out.println("withitem=" + arg0.getName());
		}

		@Override
		public void visit(ValuesStatement arg0) {
			System.out.println(arg0.getExpressions().toString());
		}
		
	}
	
	static class FilterExpressionVisitorAdapter extends ExpressionVisitorAdapter{
	    int depth = 0;
	    public void processLogicalExpression( BinaryExpression expr, String logic){
	        System.out.println(StringUtils.repeat("-", depth) + logic);

	        depth++;
	        expr.getLeftExpression().accept(this);
	        expr.getRightExpression().accept(this);
	        if(  depth != 0 ){
	            depth--;
	        }
	    }

	    @Override
	    protected void visitBinaryExpression(BinaryExpression expr) {
	        if (expr instanceof ComparisonOperator) {
	            System.out.println(StringUtils.repeat("-", depth) + 
	                "left=" + expr.getLeftExpression() + 
	                "  op=" +  expr.getStringExpression() + 
	                "  right=" + expr.getRightExpression() );
	        } 
	        super.visitBinaryExpression(expr); 
	    }

	    @Override
	    public void visit(AndExpression expr) {
	        processLogicalExpression(expr, "AND");

	    }
	    @Override
	    public void visit(OrExpression expr) {
	        processLogicalExpression(expr, "OR");
	    }
	    @Override
	    public void visit(Parenthesis parenthesis) {
	        parenthesis.getExpression().accept(this);
	    }

	}
	
	static class MyItemListVisitor implements ItemsListVisitor {

		@Override
		public void visit(SubSelect arg0) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public void visit(ExpressionList arg0) {
			
			
		}

		@Override
		public void visit(NamedExpressionList arg0) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public void visit(MultiExpressionList arg0) {
			// TODO Auto-generated method stub
			
		}
		
	}
}
